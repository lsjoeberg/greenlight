package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *application) serve() error {
	// Declare an HTTP server.
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Create a shutdownError channel to receive any errors returned by the graceful Shutdown() function.
	shutdownError := make(chan error)

	// Start a background goroutine to catch os signals.
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		app.logger.PrintInfo("shutting down server", map[string]string{
			"signal": s.String(),
		})

		// Perform a graceful shutdown after receiving interrupt or termination signals.
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Shutdown the server; only send on the shutdownError channel if it returns an error.
		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		// Wait for any background goroutines to complete their tasks.
		app.logger.PrintInfo("completing background tasks", map[string]string{"addr": srv.Addr})
		app.wg.Wait()
		shutdownError <- nil
	}()

	app.logger.PrintInfo("starting server", map[string]string{
		"addr": srv.Addr,
		"env":  app.config.env,
	})

	// Start the server. If we receive and ErrServerClosed error, it means that
	// the server has gracefully shut down, via srv.Shutdown().
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	// Otherwise, we wait to receive the return value from Shutdown() on the
	// shutdownError channel. If return value is an error, we know that there was a
	// problem with the graceful shutdown, and we return the error.
	err = <-shutdownError
	if err != nil {
		return err
	}

	// At this point we know that the graceful shutdown completed successfully, and we
	// log a "stopped server" message.
	app.logger.PrintInfo("stopped server", map[string]string{"addr": srv.Addr})

	return nil
}
