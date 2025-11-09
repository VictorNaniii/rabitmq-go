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
	"syscall"

	"github.com/rabbitmq/amqp091-go"
	grpserver "google.golang.org/grpc"
)

func main() {
	GrpcAdres := ":9093"
	rabbitMqUri := env.GetString("RABBITMQ_URI", "amq://guest:guest@localhost:5672/")

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
	conn, err := amqp091.Dial(rabbitMqUri)
	if err != nil {
		log.Fatal("Failed to connect to rabbitmq")
	}
	defer conn.Close()
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
