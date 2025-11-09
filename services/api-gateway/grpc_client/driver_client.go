package grpc_client

import (
	"os"
	pb "ride-sharing/shared/proto/driver"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type driverServiceCLient struct {
	Client pb.DriverServiceClient
	conn   *grpc.ClientConn
}

func NewDriverServiceClient() (*driverServiceCLient, error) {
	driverServiceUrl := os.Getenv("DRIVER_SERVICE_URL")
	if driverServiceUrl == "" {
		driverServiceUrl = "driver-service:9092"
	}
	conn, er := grpc.NewClient(driverServiceUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if er != nil {
		return nil, er
	}
	client := pb.NewDriverServiceClient(conn)
	return &driverServiceCLient{
		Client: client,
		conn:   conn,
	}, nil
}
func (c *driverServiceCLient) Close() {
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			return
		}
	}
}
