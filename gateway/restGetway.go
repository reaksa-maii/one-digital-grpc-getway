package gateway

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/reaksa-maii/one_digital_grpc_getway/proto/podcast/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func GetMethod() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	mux := runtime.NewServeMux()

	conn, err := grpc.NewClient("localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	defer conn.Close()

	if err := pb.RegisterPodcatServiceHandler(ctx, mux, conn); err != nil {
		return err
	}

	fmt.Println("Starting RestFull server on localhost:8081...")
	if err := http.ListenAndServe(":9090", mux); err != nil {
		return err
	}
	return nil
}
func PostMethod() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	mux := runtime.NewServeMux()

	conn, err := grpc.DialContext(ctx, ":8080", grpc.WithInsecure())
	if err != nil {
		return err
	}

	if err := pb.RegisterPodcatServiceHandler(ctx, mux, conn); err != nil {
		return err
	}

	fmt.Println("Starting Rest-Gateway server on :8081...")
	if err := http.ListenAndServe(":8081", mux); err != nil {
		return err
	}
	return nil
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
	return nil
}
