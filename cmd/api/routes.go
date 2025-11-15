package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/healthcheck", app.healthcheckHandler)
	r.Get("/posts/{id}", app.showPostHandler)
	r.Post("/posts", app.createPostHandler)
	r.Put("/posts/{id}", app.updatePostHandler)
	r.Delete("/posts/{id}", app.deletePostHandler)

	return r
}
