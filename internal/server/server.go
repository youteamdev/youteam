package server

import (
	"context"
	"errors"
	"net/http"
	"time"

	"youteam.dev/youteam/allinone/internal/config"
	"youteam.dev/youteam/allinone/internal/web"
)

// NewHTTPHandler returns the HTTP handler backed by embedded web assets.
func NewHTTPHandler() http.Handler {
	return http.FileServer(http.FS(web.FS))
}

// RunHTTPServer starts the embedded HTTP server and blocks until it exits.
func RunHTTPServer(ctx context.Context, cfg config.Config) error {
	server := &http.Server{
		Addr:              cfg.Address(),
		Handler:           NewHTTPHandler(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- server.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			return err
		}

		if err := <-errCh; err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	}
}
