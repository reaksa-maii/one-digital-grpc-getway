package utilities

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/alts"
	"log"
	"time"
)

const (
	serverAddr       = "localhost:50051"
	expectedServerSA = "maireaksa@gmail.com"
)

func authServer() {
	clientOpts := alts.DefaultClientOptions()
	clientOpts.TargetServiceAccounts = []string{expectedServerSA}
	altsTC := alts.NewClientCreds(clientOpts)
	conn, err := grpc.NewClient(serverAddr, grpc.WithTransportCredentials(altsTC))
	if err != nil {
		log.Fatalf("failed to dial server: %v", err)
	}
	defer conn.Close()

	// Example: You can now use conn to call your gRPC service
	// client := pb.NewYourServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	fmt.Println("Successfully authenticated with ALTS and connected to server.")

	// Add actual service call here if needed
	_ = ctx
}

func authToken() {

	clientOpts := alts.DefaultClientOptions()
	clientOpts.TargetServiceAccounts = []string{expectedServerSA}
	altsTC := alts.NewClientCreds(clientOpts)
	conn, err := grpc.NewClient(serverAddr, grpc.WithTransportCredentials(altsTC))
	
	if err != nil {
		log.Fatalf("failed to dial server: %v", err)
	}
	defer conn.Connect()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	fmt.Println("Successfully authenticated with ALTS and connected to server.")

	// Add actual service call here if needed
	_ = ctx
}
