package service

import (
	"context"
	"ride-sharing/services/trip-service/internal/domain"
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
