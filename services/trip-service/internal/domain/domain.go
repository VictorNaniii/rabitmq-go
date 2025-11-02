package domain

import (
	"context"
	"ride-sharing/shared/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
	trypTypes "ride-sharing/services/trip-service/pkg/types"
)

type TripModel struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID   string             `json:"userId" bson:"userId"`
	Status   string             `json:"status" bson:"status"`
	RideFare *RideFareModel
}
type TripRepository interface {
	CreateTrip(ctx context.Context, trip *TripModel) (*TripModel, error)
	SaveRideFare(ctx context.Context, f *RideFareModel) error
}
type TripService interface {
	CreateTrip(ctx context.Context, fare *RideFareModel) (*TripModel, error)
	GetRoute(ctx context.Context, pickup, destination *types.Coordinate) (*trypTypes.OSRMRoute, error)
	EstimatePackagePriceWithRoute(route *trypTypes.OSRMRoute) []*RideFareModel
	GenerateTripFares(ctx context.Context, fares []*RideFareModel, userID string) ([]*RideFareModel, error)
}
