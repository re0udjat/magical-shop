package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
}

type app struct {
	config config
	logger *slog.Logger
}

func main() {
	var cfg config

	// Flags for app
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "dev", "Envrionment (dev|stag|prod)")
	flag.Parse()

	// Structured logger which writes log entries to stdout stream
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	app := &app{
		config: cfg,
		logger: logger,
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
	err := srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}
