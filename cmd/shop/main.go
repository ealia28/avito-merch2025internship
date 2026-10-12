// Command shop запускает HTTP-сервис магазина мерча Авито.
package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/ealia28/avito-merch2025internship/internal/config"
)

func run(logger *slog.Logger) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	logger.Info("server port", "port", cfg.Server.Port)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})
	srv := &http.Server{
		Addr:              ":" + cfg.Server.Port,
		Handler:           mux,
		ReadHeaderTimeout: cfg.Server.ReadHeaderTimeout,
		ReadTimeout:       cfg.Server.ReadTimeout,
		WriteTimeout:      cfg.Server.WriteTimeout,
		IdleTimeout:       cfg.Server.IdleTimeout,
	}
	logger.Info("starting server", "port", srv.Addr)
	err = srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("server error: %w", err)
	}
	return nil
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	err := run(logger)
	if err != nil {
		logger.Error("failed to run application", "error", err)
		os.Exit(1)
	}
}
