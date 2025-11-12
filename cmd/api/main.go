package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
)

var version = "1.0.0"

type application struct {
	cfg    config
	logger *slog.Logger
}

type config struct {
	port int
	env  string
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	var cfg config

	flag.IntVar(&cfg.port, "port", 9090, "server port")
	flag.StringVar(&cfg.env, "env", "dev", "Environment (dev | prod | test)")

	flag.Parse()

	app := application{
		cfg:    cfg,
		logger: logger,
	}

	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("What's up!"))
	})

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.cfg.port),
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		ErrorLog:     slog.NewLogLogger(app.logger.Handler(), slog.LevelError),
	}

	logger.Info("server starting", "addr", srv.Addr, "env", app.cfg.env)
	if err := srv.ListenAndServe(); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
