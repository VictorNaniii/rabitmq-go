package main

import (
	"context"
	grpserver "google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"ride-sharing/services/trip-service/internal/infrastructure/grpc"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
	"ride-sharing/services/trip-service/internal/service"
	"syscall"
)

func main() {
	GrpcAdres := ":9093"

	inmemRepo := repository.NewInMemRepository()
	service := service.NewTripService(inmemRepo)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh
		cancel()
	}()
	lis, err := net.Listen("tcp", GrpcAdres)
	if err != nil {
		log.Fatal("FAiled to listen:", err)
	}
	//Starting gRPC Server
	grpcServer := grpserver.NewServer()
	grpc.NewGRPCHandler(grpcServer, service)
	log.Println("sTARTING gRPC server on port: ", lis.Addr().String())
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Println("Failed to serve:", err)
			cancel()
		}
	}()
	//wait for shudown signal
	<-ctx.Done()
	log.Println("Shutting down the server...")
	grpcServer.GracefulStop()
}
