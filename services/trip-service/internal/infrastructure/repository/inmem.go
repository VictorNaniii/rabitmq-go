package repository

import (
	"context"
	"ride-sharing/services/trip-service/internal/domain"
)

type InMemRepository struct {
	trips      map[string]*domain.TripModel
	riderFAres map[string]*domain.RideFareModel
}

func NewInMemRepository() *InMemRepository {
	return &InMemRepository{
		trips:      make(map[string]*domain.TripModel),
		riderFAres: make(map[string]*domain.RideFareModel),
	}
}

func (i *InMemRepository) CreateTrip(ctx context.Context, trip *domain.TripModel) (*domain.TripModel, error) {
	i.trips[trip.ID.Hex()] = trip
	return trip, nil
}
