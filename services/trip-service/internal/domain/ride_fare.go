package domain

import (
	pb "ride-sharing/shared/proto/trip"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RideFareModel struct {
	ID                primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserId            string             `json:"userId" bson:"userId"`
	PackageSlug       string             `json:"packageSlug" bson:"packageSlug"`
	TotalPriceInCents float64            `json:"totalPriceInCents" bson:"totalPriceInCents"`
}

func (r *RideFareModel) ToProto() *pb.RideFare {
	return &pb.RideFare{
		Id:          r.ID.Hex(),
		UserID:      r.UserId,
		PackageSLug: r.PackageSlug,
		TotalPrice:  r.TotalPriceInCents,
	}
}

func ToRideFaresProto(fares []*RideFareModel) []*pb.RideFare {
	var protoFares []*pb.RideFare
	for _, f := range fares {
		protoFares = append(protoFares, f.ToProto())
	}
	return protoFares
}
