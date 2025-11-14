package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/Infamous003/go-blog/internal/data"
	"github.com/Infamous003/go-blog/internal/validator"
)

func (app *application) showPostHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notfoundResponse(w, r)
		return
	}

	p := data.Post{
		ID:        id,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Title:     "Why I am allowed to say the N word",
		Subtitle:  "[Spoiler Alert] I'm Black, or am I?",
		Content:   "uksvhj,nwskvjh,anmsfv  aj,sdvbkjg,vals",
		Status:    "draft",
		Claps:     23,
		Version:   1,
	}

	headers := make(http.Header)
	headers.Set("Languages", "en,fr")
	err = app.writeJSON(w, http.StatusOK, envelope{"post": p}, headers)
	if err != nil {
		app.serverErrorResponse(w, r)
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
			app.serverErrorResponse(w, r)
		}
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"post": post}, nil)
	if err != nil {
		app.serverErrorResponse(w, r)
	}
}
