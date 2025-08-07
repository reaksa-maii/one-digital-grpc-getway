package gateway

import (
	"context"
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/reaksa-maii/one_digital_grpc_getway/proto/podcast/v1"
	"google.golang.org/grpc"
)

func PostMethod() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	mux := runtime.NewServeMux()

	conn, err := grpc.DialContext(ctx, ":50051", grpc.WithInsecure())
	if err != nil {
		return err
	}

	if err := pb.RegisterPodcatServiceHandler(ctx, mux, conn); err != nil {
		return err
	}

	fmt.Println("Starting Rest-Gateway server on localhost:8081...")
	if err := http.ListenAndServe(":8081", mux); err != nil {
		return err
	}
	return nil
}