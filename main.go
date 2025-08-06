package main

import (
	"github.com/reaksa-maii/one_digital_grpc_getway/gateway"
	"log"
)

func main() {

	go func() {
		if err := gateway.PostMethod(); err != nil {
			log.Fatalf("Failed to run REST server: %v", err)
		}
	}()

	if err := gateway.RungRPC(); err != nil {
		log.Fatalf("Failed to run gRPC server: %v", err)
	}
}
