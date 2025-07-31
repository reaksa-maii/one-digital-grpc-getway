package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/reaksa-maii/one_digital_grpc_getway/proto/movie/v1"
)

var ports = flag.Int("port", 50051, "port server run time")

type server struct {
	pb.UnimplementedMovieServer
}

func (s *server) UnaryMovie(_ context.Context, in *pb.MovieRequest) (*pb.MovieResponse, error) {
	fmt.Printf("Unary Movie Massage: %q\n", in.Title)
	return &pb.MovieResponse{Movie: in.Movie}, nil
}
func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *ports))
	if err != nil {
		fmt.Printf("Server not found! %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterMovieServer(s, &server{})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
