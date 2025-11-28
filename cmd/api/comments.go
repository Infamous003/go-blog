package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Infamous003/go-blog/internal/data"
	"github.com/Infamous003/go-blog/internal/validator"
	"github.com/go-chi/chi/v5"
)

func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)

	postID, err := app.readIDParam(r)
	if err != nil {
		app.notfoundResponse(w, r)
		return
	}

	_, err = app.models.Posts.Get(postID)
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
		Body string `json:"body"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	comment := data.Comment{
		Body:   input.Body,
		UserID: user.ID,
		PostID: postID,
	}

	v := validator.New()

	if data.ValidateComment(v, &comment); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Comments.Insert(&comment)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"comment": comment}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listCommentsForPostHandler(w http.ResponseWriter, r *http.Request) {
	postID, err := app.readIDParam(r)
	if err != nil {
		app.notfoundResponse(w, r)
		return
	}

	var input struct {
		Filters data.Filter
	}

	v := validator.New()
	qs := r.URL.Query()

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 2, v)

	data.ValidateFilters(v, input.Filters)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	_, err = app.models.Posts.Get(postID)
	if err != nil {
		app.notfoundResponse(w, r)
		return
	}

	comments, metadata, err := app.models.Comments.GetForPost(postID, &input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"metadata": metadata, "comments": comments}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	commentID, err := strconv.ParseInt(chi.URLParam(r, "comment_id"), 10, 0)
	if err != nil || commentID < 1 {
		app.notfoundResponse(w, r)
		return
	}

	user := app.contextGetUser(r)

	postID, err := app.readIDParam(r)
	if err != nil {
		app.notfoundResponse(w, r)
		return
	}

	err = app.models.Comments.Delete(commentID, user.ID, postID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notfoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "successfully deleted the comment"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
