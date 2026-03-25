package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/re0udjat/magic-shop/internal/data"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxConns     int
		minIdleConns int
		maxIdleTime  time.Duration
	}
}

type app struct {
	config config
	logger *slog.Logger
	models data.Models
}

func main() {
	var cfg config

	// Flags for app
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "dev", "Envrionment (dev|stag|prod)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://postgres:postgres@localhost:5432/magical_shop?sslmode=disable", "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxConns, "db-max-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.minIdleConns, "db-min-idle-conns", 5, "PostgreSQL min idle connections")
	flag.DurationVar(&cfg.db.maxIdleTime, "db-max-idle-time", 15*time.Minute, "PostgreSQL max connection idle time")
	flag.Parse()

	// Structured logger which writes log entries to stdout stream
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Open the connection pool
	dbpool, err := openDB(cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Defer a call to dbpool.Close() so that the connection pool is closed before the main()
	// function exits
	defer dbpool.Close()
	logger.Info("database connection pool established")

	app := &app{
		config: cfg,
		logger: logger,
		models: data.NewModels(dbpool),
	}

	/* Declare a server:
	*	- IdleTimeout: max duration for keep-alive connection
	*	- ReadTimeout: max duration for reading the request headers and body => Mitigate the risk from slow-client attacks
	*	- WriteTimeout: max duration for writing the response body
	 */
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	// Start the server
	logger.Info("starting server", "addr", srv.Addr, "env", cfg.env)
	err = srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}

func openDB(cfg config) (*pgxpool.Pool, error) {
	// Create a connection pool
	dbpool, err := pgxpool.New(context.Background(), cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	// Max number of opens connections in the pool (in-use + idle connections)
	dbpool.Config().MaxConns = int32(cfg.db.maxConns)

	// Min number of idle connections in the pool
	dbpool.Config().MinIdleConns = int32(cfg.db.minIdleConns)

	// Max connection idle time
	dbpool.Config().MaxConnIdleTime = cfg.db.maxIdleTime

	// Create a context with a timeout of 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Ping the database to verify the connection
	err = dbpool.Ping(ctx)
	if err != nil {
		dbpool.Close()
		return nil, err
	}

	return dbpool, nil
}
