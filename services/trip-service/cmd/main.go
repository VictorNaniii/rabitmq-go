package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	ht "ride-sharing/services/trip-service/internal/infrastructure/http"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
	"ride-sharing/services/trip-service/internal/service"
	"syscall"
	"time"
)

func main() {
	port := "8083"

	inmemRepo := repository.NewInMemRepository()
	service := service.NewTripService(inmemRepo)
	handler := ht.NewHttpControler(service)

	log.Println("Starting trip service on port", port)

	// Allow overriding port from environment
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/trips/preview", handler.HandleTripPreview)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	serverErrors := make(chan error, 1)
	shutdown := make(chan os.Signal, 1)

	// Start HTTP server in background
	go func() {
		log.Println("Server is listening on port", port)
		serverErrors <- server.ListenAndServe()
	}()

	// Catch termination signals
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Fatalf("error starting server: %v", err)

	case sig := <-shutdown:
		log.Printf("received signal %v: starting shutdown", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("graceful shutdown did not complete in %v: %v", 10*time.Second, err)
			if err := server.Close(); err != nil {
				log.Printf("error closing server: %v", err)
			}
		}
	}
}
