package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/Infamous003/go-blog/internal/validator"
)

type Comment struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    int64     `json:"user_id"`
	PostID    int64     `json:"post_id"`
	Version   int64     `json:"version"`
}

func ValidateComment(v *validator.Validator, c *Comment) {
	v.Check(c.Body != "", "comment", "must be provided")
	v.Check(len(c.Body) >= 10, "comment", "must be atleast 10 characters long")
}

type CommentModel struct {
	DB *sql.DB
}

func (m CommentModel) Insert(comment *Comment) error {
	query := `
		INSERT INTO comments (body, user_id, post_id)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, user_id, post_id, version
	`

	args := []any{
		comment.Body,
		comment.UserID,
		comment.PostID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(
		&comment.ID,
		&comment.CreatedAt,
		&comment.UserID,
		&comment.PostID,
		&comment.Version,
	)
}

func (m CommentModel) GetForPost(postID int64, filters *Filter) ([]*Comment, Metadata, error) {
	query := `
		SELECT count(*) over(), id, body, user_id, post_id, created_at, updated_at 
		FROM comments
		WHERE post_id = $1
		ORDER BY id
		LIMIT $2 OFFSET $3
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, postID, filters.limit(), filters.offset())
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	comments := []*Comment{}
	totalRecords := 0

	for rows.Next() {
		var c Comment

		err := rows.Scan(
			&totalRecords,
			&c.ID,
			&c.Body,
			&c.UserID,
			&c.PostID,
			&c.CreatedAt,
			&c.UpdatedAt,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		comments = append(comments, &c)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return comments, metadata, err
}

func (m CommentModel) Delete(commentID, userID, postID int64) error {
	query := `
		DELETE FROM comments
		WHERE id = $1 AND user_id = $2 AND post_id = $3
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	res, err := m.DB.ExecContext(ctx, query, commentID, userID, postID)
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
