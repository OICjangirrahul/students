package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/OICjangirrahul/students/internal"
	"github.com/OICjangirrahul/students/internal/config"
)

func main() {
	cfg := config.MustLoad()

	handler, err := internal.InitializeStudentHandler(cfg)
	if err != nil {
		slog.Error("failed to initialize handler", slog.String("error", err.Error()))
		os.Exit(1)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /students", handler.Create())
	mux.HandleFunc("GET /students/{id}", handler.GetByID())
	mux.HandleFunc("POST /students/login", handler.Login())

	srv := &http.Server{
		Addr:    cfg.HTTPServer.Addr,
		Handler: mux,
	}

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				slog.Error("graceful shutdown timed out.. forcing exit.")
				os.Exit(1)
			}
		}()

		// Trigger graceful shutdown
		err := srv.Shutdown(shutdownCtx)
		if err != nil {
			slog.Error("server shutdown error", slog.String("error", err.Error()))
		}
		serverStopCtx()
	}()

	// Run the server
	slog.Info("starting server", slog.String("addr", cfg.HTTPServer.Addr))
	err = srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		slog.Error("server error", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
}
