package main

import (
	"context"
	"database/sql"
	"flag"
	"log/slog"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var version = "1.0.0"

type application struct {
	cfg    config
	logger *slog.Logger
}

type config struct {
	port int
	env  string

	db struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  time.Duration
	}
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	var cfg config

	// Server configurations
	flag.IntVar(&cfg.port, "port", 9090, "server port")
	flag.StringVar(&cfg.env, "env", "dev", "Environment (dev | prod | test)")

	// DB configurations
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("GOBLOG_DSN"), "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.DurationVar(&cfg.db.maxIdleTime, "db-max-idle-time", 10*time.Minute, "PostgreSQL max connection idle time")

	flag.Parse()

	app := application{
		cfg:    cfg,
		logger: logger,
	}

	db, err := OpenDB(cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	logger.Info("database connection pool established")

	if err = app.serve(); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

// Opens a connection to the database and returns an instance of it
func OpenDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	db.SetConnMaxIdleTime(cfg.db.maxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
