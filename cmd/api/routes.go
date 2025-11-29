package main

import (
	"expvar"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()

	r.Use(app.metrics)
	r.Use(middleware.Logger)
	r.Use(app.rateLimiter)
	r.Use(app.authenticate)

	// custom error responses
	r.MethodNotAllowed(app.methodNotAllowedResponse)
	r.NotFound(app.notfoundResponse)

	r.Get("/healthcheck", app.healthcheckHandler)

	r.Handle("/metrics", expvar.Handler())

	// USERS endpoints
	r.Route("/users", func(r chi.Router) {
		r.Post("/", app.registerUserHandler)
		r.Put("/activated", app.activateUserHandler)
		r.Get("/me", app.getProfileHandler)
	})

	// TOKENS endpoints
	r.Route("/tokens", func(r chi.Router) {
		r.Post("/authentication", app.createAuthenticationTokenHandler)
	})

	// POSTS endpoints
	r.Route("/posts", func(r chi.Router) {
		r.Post("/", app.requireActivatedUser(app.createPostHandler))
		r.Get("/", app.requireActivatedUser(app.ListPostsHandler))

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", app.requireActivatedUser(app.showPostHandler))
			r.Patch("/", app.requireActivatedUser(app.updatePostHandler))
			r.Delete("/", app.requireActivatedUser(app.deletePostHandler))

			r.Post("/publish", app.requireActivatedUser(app.publishPostHandler))
			r.Post("/clap", app.requireActivatedUser(app.clapPostHandler))

			r.Route("/comments", func(r chi.Router) {
				r.Post("/", app.requireActivatedUser(app.createCommentHandler))
				r.Get("/", app.requireActivatedUser(app.listCommentsForPostHandler))
				r.Delete("/{comment_id}", app.requireActivatedUser(app.deleteCommentHandler))
				r.Patch("/{comment_id}", app.requireActivatedUser(app.updateCommentHandler))
			})
		})
	})

	return r
}
