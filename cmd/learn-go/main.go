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

	"github.com/PratikPradhan987/learn-go/internal/config"
	"github.com/PratikPradhan987/learn-go/internal/http/handlers/student"
)

func main() {
	// Load Configuration
	cfg := config.MustLoad()

	// Setup Router
	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request){
        w.Write([]byte("Hello, World!"))
	})
	router.HandleFunc("POST /api/student", student.New())

	// Server Setup

	server := http.Server {
		Addr:  cfg.Addr,
        Handler: router,
	}
	slog.Info("Server Started",slog.String("address" ,cfg.Addr))

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func(){
		err := server.ListenAndServe()

		if err!= nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	<-done
	slog.Info("Shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err!= nil {
		slog.Error("Shutting down the server",slog.String("error" ,err.Error()))
    }

	slog.Info("Server gracefully stopped")
}