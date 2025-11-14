package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

// A container representing all the models
type Models struct {
	Posts PostModel
}

// Returns a Models struct which contains all the models initialized with a DB
func NewModels(db *sql.DB) Models {
	return Models{
		Posts: PostModel{DB: db},
	}
}
