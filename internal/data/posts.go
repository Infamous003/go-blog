package data

import (
	"database/sql"
	"time"

	"github.com/Infamous003/go-blog/internal/validator"
)

type Post struct {
	ID          int64     `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Title       string    `json:"title"`
	Subtitle    string    `json:"subtitle,omitzero"`
	Content     string    `json:"content"`
	Tags        []string  `json:"tags"`
	Claps       int64     `json:"claps"`
	Status      string    `json:"status"` // Draft or Published
	PublishedAt time.Time `json:"published_at"`
	Version     int64     `json:"version"`
	Slug        string    `json:"slug"`
}

func ValidatePost(v *validator.Validator, post *Post) {
	v.Check(post.Title != "", "title", "must be provided")
	v.Check(len(post.Title) >= 25, "title", "must be atleast 25 characters long")
	v.Check(len(post.Title) <= 75, "title", "must not be longer than 100 characters")

	v.Check(len(post.Subtitle) <= 30, "subtitle", "must not be longer than 25 characters long")

	v.Check(post.Content != "", "content", "must be provided")
	v.Check(len(post.Content) <= 5000, "content", "must not be longer than 25 charcaters")

	v.Check(validator.Unique(post.Tags), "tags", "must not contain duplicate values")
	v.Check(len(post.Tags) >= 1, "tags", "must contain atleast 1 tag")
	v.Check(len(post.Tags) <= 5, "tags", "must not contain more than 5 tags")
}

// Model representing Post, which contains a DB connection
type PostModel struct {
	DB *sql.DB
}

// Inserts a Post in the DB, returns an error if failed to do so
func (m PostModel) Insert(post *Post) error {
	return nil
}

// Fetch a Post from the DB, returns an error if failed to do so
func (m PostModel) Get(id int64) (*Post, error) {
	return nil, nil
}

// Update a Post, returns an error if failed to do so
func (m PostModel) Update(post *Post) error {
	return nil
}

// Delete a Post from the DB
func (m PostModel) Delete(id int64) error {
	return nil
}
