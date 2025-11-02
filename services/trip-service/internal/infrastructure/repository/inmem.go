package repository

import (
	"context"
	"ride-sharing/services/trip-service/internal/domain"
)

type InMemRepository struct {
	trips      map[string]*domain.TripModel
	riderFares map[string]*domain.RideFareModel
}

func NewInMemRepository() *InMemRepository {
	return &InMemRepository{
		trips:      make(map[string]*domain.TripModel),
		riderFares: make(map[string]*domain.RideFareModel),
	}
}

func (i *InMemRepository) CreateTrip(ctx context.Context, trip *domain.TripModel) (*domain.TripModel, error) {
	i.trips[trip.ID.Hex()] = trip
	return trip, nil
}
func (r *InMemRepository) SaveRideFare(ctx context.Context, f *domain.RideFareModel) error {
	r.riderFares[f.ID.Hex()] = f
	return nil
}
