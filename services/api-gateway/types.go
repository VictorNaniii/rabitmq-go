package main

type previewTripRequest struct {
	UserID      string  `json:"userId"`
	PackageSlug string  `json:"packageSlug"`
	DistanceKm  float32 `json:"distanceKm"`
}
