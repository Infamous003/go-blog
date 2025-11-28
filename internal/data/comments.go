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
