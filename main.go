package main

import (
	"context"
	"github.com/zugarekd/go-track/handlers"
	"github.com/zugarekd/go-track/middleware"
	"github.com/zugarekd/go-track/server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	port := "8088"

	mux := http.NewServeMux()
	mux.HandleFunc("/log", handlers.RadonProHandler)

	wrappedMux := middleware.LoggingMiddleware(mux)

	srv := server.NewServer(wrappedMux, port)

	go func() {
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on %s: %v\n", port, err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Stop(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
