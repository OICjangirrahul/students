package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/OICjangirrahul/students/internal/config"
	"github.com/OICjangirrahul/students/internal/http/handlers/student"
	"github.com/OICjangirrahul/students/internal/storage/sqlite"
)

func main() {
	//load config
	cfg := config.MustLoad()

	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("storage initialized")

	//setup router
	router := http.NewServeMux()

	router.HandleFunc("POST /api/students", student.New(storage))
	router.HandleFunc("POST /api/students/{id}", student.GetById(storage))

	server := http.Server{
		Addr:    cfg.HTTPServer.Addr,
		Handler: router,
	}

	slog.Info("server is running", slog.String("address", cfg.HTTPServer.Addr))

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("failed to start server")
		}
	}()

	<-done

	slog.Info("shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	// err := server.Shutdown(ctx)

	// if err != nil {
	// 	slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	// }

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("Server shutdown successfully")

	// go run cmd/students/main.go -config config/local.yaml
	// go run cmd/students/main.go -config config/dev.yaml
	// go run cmd/students/main.go -config config/prod.yaml

}
