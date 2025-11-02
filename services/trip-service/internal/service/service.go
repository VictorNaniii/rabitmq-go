package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"ride-sharing/services/trip-service/internal/domain"
	trypTypes "ride-sharing/services/trip-service/pkg/types"
	"ride-sharing/shared/types"
)

type TripService struct {
	repo domain.TripRepository
}

func NewTripService(repo domain.TripRepository) *TripService {
	return &TripService{repo: repo}
}

func (s *TripService) CreateTrip(ctx context.Context, fare *domain.RideFareModel) (*domain.TripModel, error) {
	trip := &domain.TripModel{
		UserID:   fare.UserId,
		Status:   "created",
		RideFare: fare,
	}
	return s.repo.CreateTrip(ctx, trip)
}
func (s *TripService) GetRoute(ctx context.Context, pickup, destination *types.Coordinate) (*trypTypes.OSRMRoute, error) {

	url := fmt.Sprintf(
		"http://router.project-osrm.org/route/v1/driving/%f,%f;%f,%f?overview=full&geometries=geojson",
		pickup.Longitude, pickup.Latitude, destination.Longitude, destination.Latitude,
	)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error making request to OSRM:", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	var routeResponse trypTypes.OSRMRoute

	if err := json.Unmarshal(body, &routeResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal OSRM response: %w", err)
	}
	return &routeResponse, nil
}
