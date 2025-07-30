package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/reaksa-maii/one_digital_grpc_getway/proto/echo/v2"
)

var ports = flag.Int("port", 50052, "port server run time")

type EchoService struct {
	pb.UnimplementedEchoServer
}

func (s *EchoService) EchoMessage(ctx context.Context, req *pb.EchoRequest) (*pb.EchoReply, error) {
	return &pb.EchoReply{Message: "Hello gRPC"}, nil
}

func main() {

	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *ports))
	if err != nil {
		log.Fatalf("Port not work %v", err)
	}
	fmt.Printf("Server list %v", lis.Addr())
	s := grpc.NewServer()
	pb.RegisterEchoServer(s, &EchoService{})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
