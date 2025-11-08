package domain

import (
	"context"
	trypTypes "ride-sharing/services/trip-service/pkg/types"

	"go.mongodb.org/mongo-driver/bson/primitive"

	pb "ride-sharing/shared/proto/trip"
	"ride-sharing/shared/types"
)

type TripModel struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID   string             `json:"userId" bson:"userId"`
	Status   string             `json:"status" bson:"status"`
	RideFare *RideFareModel
	Driver   *pb.TripDriver
}
type TripRepository interface {
	CreateTrip(ctx context.Context, trip *TripModel) (*TripModel, error)
	SaveRideFare(ctx context.Context, f *RideFareModel) error
	GetRiderFarerByID(ctx context.Context, id string) (*RideFareModel, error)
}
type TripService interface {
	CreateTrip(ctx context.Context, fare *RideFareModel) (*TripModel, error)
	GetRoute(ctx context.Context, pickup, destination *types.Coordinate) (*trypTypes.OSRMRoute, error)
	EstimatePackagesPriceWithRoute(route *trypTypes.OSRMRoute) []*RideFareModel
	GenerateTripFares(ctx context.Context, fares []*RideFareModel, userID string) ([]*RideFareModel, error)
	GetAndValidateFare(ctx context.Context, fareID, userID string) (*RideFareModel, error)
}
