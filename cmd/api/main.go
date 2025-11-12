package main

import (
	"flag"
	"log/slog"
	"os"
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

	if err := app.serve(); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
