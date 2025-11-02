package grpc

import (
	"context"
	"log"
	"ride-sharing/services/trip-service/internal/domain"
	pb "ride-sharing/shared/proto/trip"
	"ride-sharing/shared/types"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type gRPCHandler struct {
	pb.UnimplementedTripServiceServer
	service domain.TripService
}

func NewGRPCHandler(server *grpc.Server, service domain.TripService) *gRPCHandler {
	handler := &gRPCHandler{
		service: service,
	}
	pb.RegisterTripServiceServer(server, handler)
	return handler
}
func (h *gRPCHandler) CreateTrip(ctx context.Context, req *pb.CreateTripRequest) (*pb.CreateTripResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateTrip not implemented")
}

func (h *gRPCHandler) PreviewTrip(ctx context.Context, request *pb.PreviewTripRequest) (*pb.PreviewTripResponse, error) {
	pickup := request.GetStartLocation()
	destination := request.GetEndLocation()
	pickupCord := &types.Coordinate{
		Latitude:  pickup.Latitude,
		Longitude: pickup.Longitude,
	}
	destinationCord := &types.Coordinate{
		Latitude:  destination.Latitude,
		Longitude: destination.Longitude,
	}
	t, err := h.service.GetRoute(ctx, pickupCord, destinationCord)
	if err != nil {
		log.Print(err)
		return nil, status.Errorf(codes.Internal, "Failed to get route rute: %v", err)
	}
	dataRoute, err := t.ToProto()
	if err != nil {
		log.Println(err)
	}
	return &pb.PreviewTripResponse{
		Route:     dataRoute,
		RideFares: []*pb.RideFare{},
	}, nil
}
