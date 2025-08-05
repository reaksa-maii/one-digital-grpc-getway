package main

import (
	"log"
	"github.com/reaksa-maii/one_digital_grpc_getway/utilities"
)

func main() {
	go func() {
		if err := utilities.GetMethod(); err != nil {
			log.Fatalf("Failed to run REST server: %v", err)
		}
	}()

	if err := utilities.RungRPC(); err != nil {
		log.Fatalf("Failed to run gRPC server: %v", err)
	}
}
