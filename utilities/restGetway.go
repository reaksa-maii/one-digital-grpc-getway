package utilities

import (
	"context"
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/reaksa-maii/one_digital_grpc_getway/proto/movie/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/alts"
)

const (
	serverAdd         = "localhost:50051"
	expectedServerSAC = "maireaksa@gmail.com"
)

func getMethod() error {
	ctx := context.Background()
	mux := runtime.NewServeMux()

	err := grpc.DialOption(ctx, "localhost:8080", grpc.WithInsecure())
	if err != nil {
		return err
	}
	fmt.Println("Starting gRPC-Gateway server on :8081...")
	if err := http.ListenAndServe(":8081", mux); err != nil {
		return err
	}
	return nil
}
func postMethod() error {
	ctx := context.Background()
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithInsecure()} // or use credentials.NewClientTLSFromCert() for secure connection

	err := pb.RegisterMovieHandlerFromEndpoint(ctx, mux, "localhost:8080", opts)
	if err != nil {
		return err
	}

	fmt.Println("Starting gRPC-Gateway server on :8081...")
	return http.ListenAndServe(":8081", mux)
}
