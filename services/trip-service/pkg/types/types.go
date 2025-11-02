package types

import pb "ride-sharing/shared/proto/trip"

type OSRMRoute struct {
	Route []struct {
		Distance float64 `json:"distance"`
		Duration float64 `json:"duration"`
		Geometry struct {
			Cordinates [][]float64 `json:"coordinates"`
		} `json:"geometry"`
	} `json:"routes"`
}

func (o *OSRMRoute) ToProto() (*pb.Route, error) {
	route := o.Route[0]
	geometry := route.Geometry.Cordinates
	cordinates := make([]*pb.Coordinate, len(geometry))
	for i, cord := range geometry {
		cordinates[i] = &pb.Coordinate{
			Longitude: cord[0],
			Latitude:  cord[1],
		}
	}
	return &pb.Route{
		Geometry: []*pb.Geometry{
			{
				Coordinates: cordinates,
			},
		},
		Distance: route.Distance,
		Duration: route.Duration,
	}, nil
}
