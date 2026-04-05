package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *app) serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(app.logger.Handler(), slog.LevelError),
	}

	// Receive any errors returned by the graceful Shutdown() function
	shutdownError := make(chan error)

	go func() {
		// Create a quit channel which carries os.Signal
		quit := make(chan os.Signal, 1)

		// Listen for incoming SIGINT and SIGTERM signals and relay them to the quit
		// channel
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		// Block until we receive a signal from the quit channel
		s := <-quit
		app.logger.Info("shutting down server", "signal", s.String())

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Only send Shutdown() on the shutdownError channel if it returns an error
		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		app.logger.Info("completing background tasks", "addr", srv.Addr)

		// Blocking until the backgroun goroutines have finished. Then return nil on the
		// shutdownError channel, to indicate that the shutdown complted without any issues
		app.wg.Wait()
		shutdownError <- nil
	}()

	app.logger.Info("starting server", "addr", srv.Addr, "env", app.config.env)

	// If we see this error, it's actually a good thing and an indication that the graceful
	// shutdown has started.
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	// Wait to receive the return value from Shutdown() on the shutdownError channel
	err = <-shutdownError
	if err != nil {
		return err
	}

	app.logger.Info("stopped server", "addr", srv.Addr)
	return nil
}
