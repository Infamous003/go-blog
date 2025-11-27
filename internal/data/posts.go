package data

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/Infamous003/go-blog/internal/validator"
	"github.com/lib/pq"
)

type Post struct {
	ID          int64      `json:"id"`
	CreatedAt   time.Time  `json:"created_at,omitzero"`
	UpdatedAt   time.Time  `json:"updated_at,omitzero"`
	Title       string     `json:"title"`
	Subtitle    string     `json:"subtitle,omitzero"`
	Content     string     `json:"content"`
	Tags        []string   `json:"tags"`
	Claps       int64      `json:"claps"`
	Status      string     `json:"status,omitzero"` // Draft or Published
	PublishedAt *time.Time `json:"published_at"`    // when it in null in the db, json response automatically fills the time as 0.000, and you don't want that, so keep it a pointer
	Version     int64      `json:"version,omitzero"`
	Slug        string     `json:"slug"`
}

func ValidatePost(v *validator.Validator, post *Post) {
	v.Check(post.Title != "", "title", "must be provided")
	v.Check(len(post.Title) >= 10, "title", "must be atleast 25 characters long")
	v.Check(len(post.Title) <= 120, "title", "must not be longer than 100 characters")

	v.Check(len(post.Subtitle) <= 200, "subtitle", "must not be longer than 25 characters long")

	v.Check(post.Content != "", "content", "must be provided")
	v.Check(len(post.Content) >= 20, "content", "must be provided")
	v.Check(len(post.Content) <= 10000, "content", "must not be longer than 25 charcaters")

	v.Check(validator.Unique(post.Tags), "tags", "must not contain duplicate values")
	v.Check(len(post.Tags) >= 1, "tags", "must contain atleast 1 tag")
	v.Check(len(post.Tags) <= 5, "tags", "must not contain more than 5 tags")
}

func (p *Post) GenerateSlug() {
	// can use Split as well. Fields handles multiple spaces properly
	words := strings.Fields(strings.ToLower(p.Title))
	p.Slug = strings.Join(words, "-")
}

// Model representing Post, which contains a DB connection
type PostModel struct {
	DB *sql.DB
}

// Inserts a Post in the DB, returns an error if failed to do so
func (m PostModel) Insert(post *Post) error {
	query := `
		INSERT INTO posts (title, subtitle, content, tags, slug)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at, slug, claps, status, version
	`
	args := []any{post.Title, post.Subtitle, post.Content, pq.Array(post.Tags), post.Slug}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.Slug,
		&post.Claps,
		&post.Status,
		&post.Version,
	)

	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "posts_slug_key"`:
			return ErrDuplicateSlug
		default:
			return err
		}
	}

	return nil
}

// Fetch a Post from the DB, returns an error if failed to do so
func (m PostModel) Get(id int64) (*Post, error) {
	query := `
		SELECT id, created_at, title, subtitle, content, tags, status, claps, slug, updated_at, published_at, version
		FROM posts
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var post Post

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.Title,
		&post.Subtitle,
		&post.Content,
		pq.Array(&post.Tags),
		&post.Status,
		&post.Claps,
		&post.Slug,
		&post.UpdatedAt,
		&post.PublishedAt,
		&post.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &post, nil
}

func (m PostModel) GetAll(title string, tags []string, filters Filter) ([]*Post, Metadata, error) {
	query := `
		SELECT 
			count(*) OVER(),
			id,
			slug,
			title,
			subtitle,
			published_at,
			tags,
			claps
		FROM posts
		WHERE status = 'published' 
			AND (search_vector @@ plainto_tsquery('english', $1) OR $1 = '')
			AND (tags @> $2 OR $2 = '{}') 
		ORDER BY
			(CASE WHEN $1 = '' THEN 1 ELSE 0 END),  -- 0 = search mode, 1 = normal mode
			ts_rank(search_vector, plainto_tsquery('english', $1)) DESC,
			published_at DESC
		LIMIT $3 OFFSET $4
	`

	/*
		Search behavior:
		1. plainto_tsquery() turns the raw search string into a tsquery.
		2. search_vector @@ tsquery matches posts based on weighted full-text search
		(title = A, subtitle = B, content = C).
		3. If the search string ($1) is empty, the full-text search filter is skipped.

		Ordering behavior:
		We use a CASE expression to switch between two sorting modes:

		- When $1 = ''  → normal listing mode
			* CASE returns 1
			* ts_rank is ignored (all rows have rank 0 anyway)
			* posts are ordered by published_at DESC

		- When $1 != '' → search mode
			* CASE returns 0
			* rows are ordered primarily by ts_rank DESC (relevance)
			* ties are broken using published_at DESC

		This avoids mixing incompatible types in a CASE expression (timestamp vs real)
		and cleanly falls back to date sorting when no search term is provided.
	*/

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{title, pq.Array(tags), filters.limit(), filters.offset()}

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	posts := []*Post{}
	totalRecords := 0

	for rows.Next() {
		var post Post

		err := rows.Scan(
			&totalRecords,
			&post.ID,
			&post.Slug,
			&post.Title,
			&post.Subtitle,
			&post.PublishedAt,
			pq.Array(&post.Tags),
			&post.Claps,
		)

		if err != nil {
			return nil, Metadata{}, err
		}

		posts = append(posts, &post)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return posts, metadata, nil
}

// Update a Post, returns an error if failed to do so
func (m PostModel) Update(post *Post) error {
	query := `
		UPDATE posts
		SET title = $1,
			subtitle = $2, 
			content = $3, 
			tags = $4, 
			slug = $5, 
			version = version + 1, 
			updated_at = NOW()
		WHERE 
			id = $6 AND version = $7
		RETURNING version
	`
	args := []any{post.Title, post.Subtitle, post.Content, pq.Array(post.Tags), post.Slug, post.ID, post.Version}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&post.Version)
	if err != nil {
		switch {
		// if the version was changed, then you wont find the exact row, which means it was edited
		// hence an edit conflict
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		case err.Error() == `pq: duplicate key value violates unique constraint "posts_slug_key"`:
			return ErrDuplicateSlug
		default:
			return err
		}
	}

	return nil
}

// Delete a Post from the DB
func (m PostModel) Delete(id int64, userID int64) error {
	query := `
		DELETE FROM posts
		WHERE id = $1 AND user_id = $2
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	res, err := m.DB.ExecContext(ctx, query, id, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

// Publishes a post, updating its status and published_at field
func (m PostModel) Publish(post *Post) error {
	query := `
		UPDATE posts
		SET status = 'published',
			published_at = NOW(),
			version = version + 1
		WHERE id = $1 AND version = $2
		RETURNING status, published_at, version
	`
	args := []any{post.ID, post.Version}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&post.Status, &post.PublishedAt, &post.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrRecordNotFound
		default:
			return err
		}
	}

	return nil
}

func (m PostModel) IncrementClap(id int64) error {
	query := `
		UPDATE posts
		SET claps = claps + 1
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
