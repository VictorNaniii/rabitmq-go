package grpc

import (
	"context"
	"log"
	"ride-sharing/services/trip-service/internal/domain"
	"ride-sharing/services/trip-service/internal/infrastructure/events"
	pb "ride-sharing/shared/proto/trip"
	"ride-sharing/shared/types"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type gRPCHandler struct {
	pb.UnimplementedTripServiceServer
	service   domain.TripService
	publisher *events.TripEventPublisher
}

func NewGRPCHandler(server *grpc.Server, service domain.TripService, publisher *events.TripEventPublisher) *gRPCHandler {
	handler := &gRPCHandler{
		service:   service,
		publisher: publisher,
	}
	pb.RegisterTripServiceServer(server, handler)
	return handler
}
func (h *gRPCHandler) CreateTrip(ctx context.Context, req *pb.CreateTripRequest) (*pb.CreateTripResponse, error) {
	fareID := req.GetRideFareID()
	userID := req.GetRideFareID()

	riderFare, err := h.service.GetAndValidateFare(ctx, fareID, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to validate the fare:%v ", err)
	}
	trip, err := h.service.CreateTrip(ctx, riderFare)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create a trip:%v", err)
	}

	//
	err = h.publisher.PublishTripCreate(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create the trip: %v", err)
	}
	//
	return &pb.CreateTripResponse{
		TripID: trip.ID.Hex(),
	}, nil
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

	userID := request.GetUserID()

	route, err := h.service.GetRoute(ctx, pickupCord, destinationCord, true)
	if err != nil {
		log.Print(err)
		return nil, status.Errorf(codes.Internal, "Failed to get route rute: %v", err)
	}
	estimatedFares := h.service.EstimatePackagesPriceWithRoute(route)
	fares, err := h.service.GenerateTripFares(ctx, estimatedFares, userID, route)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to generate the rider fares", err)
	}
	dataRoute, err := route.ToProto()
	if err != nil {
		log.Println(err)
	}
	return &pb.PreviewTripResponse{
		Route:     dataRoute,
		RideFares: domain.ToRideFaresProto(fares),
	}, nil
}
