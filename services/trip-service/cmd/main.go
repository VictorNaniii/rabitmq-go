package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"ride-sharing/services/trip-service/internal/infrastructure/grpc"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
	"ride-sharing/services/trip-service/internal/service"
	"ride-sharing/shared/env"
	"ride-sharing/shared/messaging"
	"syscall"

	grpserver "google.golang.org/grpc"
)

func main() {
	GrpcAdres := ":9093"
	rabbitMqUri := env.GetString("RABBITMQ_URI", "amqp://guest:guest@localhost:5672/")

	inmemRepo := repository.NewInMemRepository()
	tripSvc := service.NewTripService(inmemRepo)
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
	rabitMq, err := messaging.NewRabbitMq(rabbitMqUri)
	if err != nil {
		log.Fatalf("failed to connect rabbitmq: %v", err)
	}
	defer rabitMq.Close()
	//Starting gRPC Server
	grpcServer := grpserver.NewServer()
	grpc.NewGRPCHandler(grpcServer, tripSvc)
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
