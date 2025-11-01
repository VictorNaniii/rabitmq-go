package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"ride-sharing/shared/env"
	"time"
)

var (
	httpAddr = env.GetString("HTTP_ADDR", ":8081")
)

func main() {
	log.Println("Starting API Gateway")
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello from API Gateway"))
	})
	mux.HandleFunc("/trips/preview", handleTripPreview)

	server := &http.Server{
		Addr:    httpAddr,
		Handler: mux,
	}
	serverErrors := make(chan error, 1)
	go func() {
		log.Println("Server is listening on", httpAddr)
		serverErrors <- server.ListenAndServe()
	}()
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)

	select {
	case err := <-serverErrors:
		log.Fatalf("error starting server: %v", err)
	case sig := <-shutdown:
		log.Printf("received signal %v: starting shutdown", sig)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("graceful shutdown did not complete in %v: %v", 10*time.Second, err)
			server.Close()
		}
	}

}
