package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/reaksa-maii/one_digital_grpc_getway/proto/movie/v3"
	"github.com/reaksa-maii/one_digital_grpc_getway/utilities"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var port = flag.Int("port", 50051, "server start")

type server struct {
	pb.UnimplementedMovieServer
}

func (s *server) UnaryMovie(ctx context.Context, in *pb.MovieRequest) (*pb.MovieResponse, error) {
	fmt.Printf("Movie Unary Service")
	return &pb.MovieResponse{MovieSize: in.MovieSize}, nil
}
func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(grpc.UnaryInterceptor(utilities.UnaryInterception), grpc.StreamInterceptor(utilities.StreamingInterceptor))
	pb.RegisterMovieServer(s, &server{})
	reflection.Register(s)
	fmt.Printf("server listening at %v\n", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
