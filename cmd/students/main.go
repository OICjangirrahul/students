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
)



func main() {
	//load config
	cfg := config.MustLoad()

	

	//setup router
	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		//convert string into byte
		w.Write([]byte("Welcome to students api"))
	})


	server := http.Server{
		Addr: cfg.Addr,
		Handler: router,
	}


	slog.Info("server is running", slog.String("address", cfg.Addr))

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func () {
		err :=  server.ListenAndServe()
		if err != nil {
			log.Fatal("failed to start server")
		}
	}()

	<- done

	slog.Info("shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)

	defer cancel()

	// err := server.Shutdown(ctx)

	// if err != nil {
	// 	slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	// }


	if err := server.Shutdown(ctx); err != nil{
		slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("Server shutdown successfully")




	// go run cmd/students/main.go -config config/local.yaml
	// go run cmd/students/main.go -config config/dev.yaml
	// go run cmd/students/main.go -config config/prod.yaml

}