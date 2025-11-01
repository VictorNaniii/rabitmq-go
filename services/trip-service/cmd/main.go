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
	services := service.NewTripService(inmemRepo)
	hadler := ht.NewHttpControler(services)
	//service.NewTripService()
	log.Println("Starting trip service on port", port)
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}
	mux := http.NewServeMux()

	mux.HandleFunc("/trips/preview", hadler.HandleTripPreview)
	server := &http.Server{
		Addr:    ":8083",
		Handler: mux,
	}
	serverErrors := make(chan error, 1)
	go func() {
		log.Println("Server is listening on port", port)
		serverErrors <- server.ListenAndServe()
	}()
	shotDown := make(chan os.Signal, 1)

	signal.Notify(shotDown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Fatalf("error starting server: %v", err)
	case sig := <-shotDown:
		log.Printf("received signal %v: starting shutdown", sig)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("graceful shutdown did not complete in %v: %v", 10*time.Second, err)
		}
	}
}
