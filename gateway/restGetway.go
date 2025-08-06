package gateway

import (
	"context"
	"fmt"
	"net/http"
	"log"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/reaksa-maii/one_digital_grpc_getway/proto/podcast/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func GetMethod() error {
	ctx := context.Background()
	mux := runtime.NewServeMux()

	conn, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()

	if err := pb.RegisterPodcatServiceHandler(ctx, mux, conn); err != nil {
		return err
	}

	fmt.Println("Starting gRPC-Gateway server on localhost:8081...")
	if err := http.ListenAndServe("localhost:8081", mux); err != nil {
		return err
	}
	return nil
}
func PostMethod() error {
	ctx := context.Background()
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithInsecure()} // or use credentials.NewClientTLSFromCert() for secure connection

	err := pb.RegisterPodcatServiceHandlerFromEndpoint(ctx, mux, "localhost:8080", opts)
	if err != nil {
		return err
	}

	fmt.Println("Starting gRPC-Gateway server on localhost:8081...")
	return http.ListenAndServe("localhost:8081", mux)
}

func RunRestFull() error {
	go func() {
		if err := GetMethod(); err != nil {
			log.Fatalf("Failed to run REST server: %v", err)
		}
	}()

	if err := PostMethod(); err != nil {
		log.Fatalf("Failed to run gRPC server: %v", err)
	}
	return nil;
}
