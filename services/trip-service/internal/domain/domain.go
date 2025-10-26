package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TripModel struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID   string             `json:"userId" bson:"userId"`
	Status   string             `json:"status" bson:"status"`
	RideFare *RideFareModel
}
type TripRepository interface {
	CreateTrip(ctx context.Context, trip *TripModel) (*TripModel, error)
}
type TripService interface {
	CreateTrip(ctx context.Context, fare *RideFareModel) (*TripModel, error)
}
