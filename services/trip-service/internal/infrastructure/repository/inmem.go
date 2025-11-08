package repository

import (
	"context"
	"fmt"
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

func (r *InMemRepository) CreateTrip(ctx context.Context, trip *domain.TripModel) (*domain.TripModel, error) {
	r.trips[trip.ID.Hex()] = trip
	return trip, nil
}
func (r *InMemRepository) SaveRideFare(ctx context.Context, f *domain.RideFareModel) error {
	r.riderFares[f.ID.Hex()] = f
	return nil
}
func (r *InMemRepository) GetRiderFarerByID(ctx context.Context, id string) (*domain.RideFareModel, error) {
	fare, exist := r.riderFares[id]
	if !exist {
		return nil, fmt.Errorf("fare does not exist with ID: %s", id)
	}
	return fare, nil
}
