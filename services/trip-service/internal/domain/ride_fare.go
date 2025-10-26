package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type RideFareModel struct {
	ID                primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserId            string             `json:"userId" bson:"userId"`
	PackageSlug       string             `json:"packageSlug" bson:"packageSlug"`
	TotalPriceInCents float32            `json:"totalPriceInCents" bson:"totalPriceInCents"`
}
