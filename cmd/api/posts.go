package main

import (
	"errors"
	"net/http"

	"github.com/Infamous003/go-blog/internal/data"
	"github.com/Infamous003/go-blog/internal/validator"
)

func (app *application) showPostHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notfoundResponse(w, r)
		return
	}

	post, err := app.models.Posts.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notfoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"post": post}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title    string   `json:"title"`
		Subtitle string   `json:"subtitle"`
		Content  string   `json:"content"`
		Tags     []string `json:"tags"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	post := &data.Post{
		Title:    input.Title,
		Subtitle: input.Subtitle,
		Tags:     input.Tags,
		Content:  input.Content,
	}

	v := validator.New()

	data.ValidatePost(v, post)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	post.GenerateSlug()

	err = app.models.Posts.Insert(post)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateSlug):
			v.AddError("slug", "a post with this slug already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.logger.Error(err.Error())
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"post": post}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notfoundResponse(w, r)
		return
	}

	post, err := app.models.Posts.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notfoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Title    *string  `json:"title"`
		Subtitle *string  `json:"subtitle"`
		Content  *string  `json:"content"`
		Tags     []string `json:"tags"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	oldTitle := post.Title

	// we are only updating the values that are not nil, the ones user provided
	if input.Title != nil {
		post.Title = *input.Title
	}
	if input.Subtitle != nil {
		post.Subtitle = *input.Subtitle
	}
	if input.Content != nil {
		post.Content = *input.Content
	}
	if input.Tags != nil {
		post.Tags = input.Tags
	}

	if input.Title != nil && post.Title != oldTitle {
		post.GenerateSlug() // only generating the slug, if the title was changed
	}

	v := validator.New()

	data.ValidatePost(v, post)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Posts.Update(post)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		case errors.Is(err, data.ErrDuplicateSlug):
			v.AddError("slug", "a post with this slug already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"post": post}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notfoundResponse(w, r)
		return
	}

	err = app.models.Posts.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notfoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "post successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
