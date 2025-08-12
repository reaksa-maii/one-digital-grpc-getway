package main

import (
	"context"
	"flag"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/reaksa-maii/one_digital_grpc_getway/gateway"
	pb "github.com/reaksa-maii/one_digital_grpc_getway/proto/podcast/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
)

func main() {
	var (
		grpcAddr = flag.String("grpc", ":50051", "gRPC server listen address")
		httpAddr = flag.String("http", ":8081", "HTTP (grpc-gateway) listen address")
		restAPI  = flag.String("rest", "https://power-stone.testing.sabay.com", "Upstream REST base URL")
	)
	flag.Parse()

	// start gRPC server
	lis, err := net.Listen("tcp", *grpcAddr)
	if err != nil {
		log.Fatalf("listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	authorSrv := gateway.NewAuthorServer(*restAPI)
	pb.RegisterAuthorServiceServer(grpcServer, authorSrv)
	reflection.Register(grpcServer)
	go func() {
		log.Printf("gRPC server listening at %s", *grpcAddr)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("gRPC server Serve: %v", err)
		}
	}()

	// start HTTP gateway
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()} // using insecure because gateway will dial local gRPC
	// if you want TLS between gateway and gRPC, configure accordingly.

	err = pb.RegisterAuthorServiceHandlerFromEndpoint(ctx, mux, *grpcAddr, opts)
	if err != nil {
		log.Fatalf("failed to register gateway: %v", err)
	}

	log.Printf("HTTP gateway listening at %s", *httpAddr)
	if err := http.ListenAndServe(*httpAddr, mux); err != nil {
		log.Fatalf("HTTP server error: %v", err)
	}
}
