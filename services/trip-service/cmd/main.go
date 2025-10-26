package cmd

import (
	"log"
	"net/http"
	"os"
	ht "ride-sharing/services/trip-service/internal/infrastructure/http"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
	"ride-sharing/services/trip-service/internal/service"
)

func main() {
	port := "8083"
	inmemRepo := repository.NewInMemRepository()
	services := service.NewTripService(inmemRepo)
	hadler := ht.NewHttpControler(services)
	//service.NewTripService()
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}
	mux := http.NewServeMux()

	mux.HandleFunc("POST /trips/preview", hadler.HandleTripPreview)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
