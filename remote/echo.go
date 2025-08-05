package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/reaksa-maii/one_digital_grpc_getway/proto/echo/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

// Define the server struct, which implements the pb.UnimplementedGreeterServiceServer interface
type server struct {
	pb.UnimplementedGreeterServiceServer
}

// Implement the SayHello method of the pb.GreeterServiceServer interface
func (s *server) SayHello(ctx context.Context, req *pb.GreeterServiceSayHelloRequest) (*pb.GreeterServiceSayHelloResponse, error) {
	return &pb.GreeterServiceSayHelloResponse{
		Message: fmt.Sprintf("Hello, %s!", req.Name),
	}, nil
}

// Set up the gRPC server on port 8080 and serve requests indefinitely
func runGRPCServer() error {
	lis, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	pb.RegisterGreeterServiceServer(s, &server{})

	// Enable reflection to allow clients to query the server's services
	reflection.Register(s)

	fmt.Println("Starting gRPC server on localhost:8080...")
	if err := s.Serve(lis); err != nil {
		return err
	}

	return nil
}

// Set up the REST server on port 8081 and handle requests by proxying them to the gRPC server
func runRESTServer() error {
	ctx := context.Background()
	mux := runtime.NewServeMux()

	conn, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()

	if err := pb.RegisterGreeterServiceHandler(ctx, mux, conn); err != nil {
		return err
	}

	fmt.Println("Starting gRPC-Gateway server on localhost:8081...")
	if err := http.ListenAndServe("localhost:8081", mux); err != nil {
		return err
	}
	return nil
}
func main() {
	go func() {
		if err := runRESTServer(); err != nil {
			log.Fatalf("Failed to run REST server: %v", err)
		}
	}()

	if err := runGRPCServer(); err != nil {
		log.Fatalf("Failed to run gRPC server: %v", err)
	}
}
