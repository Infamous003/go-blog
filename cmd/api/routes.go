package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(app.rateLimiter)

	// custom error responses
	r.MethodNotAllowed(app.methodNotAllowedResponse)
	r.NotFound(app.notfoundResponse)

	r.Get("/healthcheck", app.healthcheckHandler)

	r.Get("/posts/{id}", app.showPostHandler)
	r.Get("/posts", app.ListPostsHandler)
	r.Post("/posts", app.createPostHandler)
	r.Post("/posts/{id}/publish", app.publishPostHandler)
	r.Patch("/posts/{id}", app.updatePostHandler)
	r.Delete("/posts/{id}", app.deletePostHandler)
	r.Post("/posts/{id}/clap", app.clapPostHandler)

	// user endpoints
	r.Post("/users", app.registerUserHandler)

	r.Put("/users/activated", app.activateUserHandler)

	return r
}
