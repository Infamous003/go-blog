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
	r.Use(app.authenticate)

	// custom error responses
	r.MethodNotAllowed(app.methodNotAllowedResponse)
	r.NotFound(app.notfoundResponse)

	r.Get("/healthcheck", app.healthcheckHandler)

	r.Get("/posts/{id}", app.requireActivatedUser(app.showPostHandler))
	r.Get("/posts", app.requireActivatedUser(app.ListPostsHandler))

	r.Post("/posts", app.requireActivatedUser(app.createPostHandler))
	r.Post("/posts/{id}/publish", app.requireActivatedUser(app.publishPostHandler))
	r.Post("/posts/{id}/clap", app.requireActivatedUser(app.clapPostHandler))

	r.Patch("/posts/{id}", app.requireActivatedUser(app.updatePostHandler))
	r.Delete("/posts/{id}", app.requireActivatedUser(app.deletePostHandler))

	// user endpoints
	r.Post("/users", app.registerUserHandler)
	r.Get("/users/me", app.getProfileHandler)

	r.Put("/users/activated", app.activateUserHandler)

	r.Post("/tokens/authentication", app.createAuthenticationTokenHandler)

	r.Post("/posts/{id}/comments", app.requireActivatedUser(app.createCommentHandler))

	return r
}
