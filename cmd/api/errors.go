package main

import "net/http"

func (app *application) logError(r *http.Request, err error) {
	app.logger.Error(err.Error(), "method", r.Method, "uri", r.URL.RequestURI())
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message string) {
	e := envelope{
		"error": message,
	}
	err := app.writeJSON(w, status, e, nil)
	if err != nil {
		app.logError(r, err)
	}
}

func (app *application) notfoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	app.errorResponse(w, r, http.StatusNotFound, message)
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request) {
	message := "the server encountered an error and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}
