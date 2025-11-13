package main

import (
	"net/http"
	"time"
)

func (app *application) showPostHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notfoundResponse(w, r)
		return
	}

	type post struct {
		ID        int64     `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		Title     string    `json:"title"`
		Subtitle  string    `json:"subtitle"`
		Content   string    `json:"content"`
		Claps     int64     `json:"claps"`
	}

	p := post{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Why I am allowed to say the N word",
		Subtitle:  "[Spoiler Alert] I'm Black, or am I?",
		Content:   "uksvhj,nwskvjh,anmsfv  aj,sdvbkjg,vals",
		Claps:     23,
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
		Title    string `json:"title"`
		Subtitle string `json:"subtitle"`
		Content  string `json:"content"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"post": input}, nil)
	if err != nil {
		app.serverErrorResponse(w, r)
	}
}
