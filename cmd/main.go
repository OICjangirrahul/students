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

	handlers, err := internal.InitializeHandlers(cfg)
	if err != nil {
		slog.Error("failed to initialize handlers", slog.String("error", err.Error()))
		os.Exit(1)
	}

	mux := http.NewServeMux()

	// Student routes
	mux.HandleFunc("POST /students", handlers.Student.Create())
	mux.HandleFunc("GET /students/{id}", handlers.Student.GetByID())
	mux.HandleFunc("POST /students/login", handlers.Student.Login())

	// Teacher routes
	mux.HandleFunc("POST /teachers", handlers.Teacher.Create())
	mux.HandleFunc("GET /teachers/{id}", handlers.Teacher.GetByID())
	mux.HandleFunc("PUT /teachers/{id}", handlers.Teacher.Update())
	mux.HandleFunc("DELETE /teachers/{id}", handlers.Teacher.Delete())
	mux.HandleFunc("POST /teachers/login", handlers.Teacher.Login())
	mux.HandleFunc("POST /teachers/{teacherId}/students/{studentId}", handlers.Teacher.AssignStudent())
	mux.HandleFunc("GET /teachers/{teacherId}/students", handlers.Teacher.GetStudents())

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
