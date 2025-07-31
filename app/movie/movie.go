package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	
	pb "github.com/reaksa-maii/one_digital_grpc_getway/proto/movie/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var ports = flag.Int("port", 50051, "port server run time")

type server struct {
	pb.UnimplementedMovieServer
}

func (s *server) UnaryMovie(_ context.Context, in *pb.MovieRequest) (*pb.MovieResponse, error) {
	fmt.Printf("Unary Movie Massage: %q\n", in.Title)
	return &pb.MovieResponse{Movie: in.Movie}, nil
}
func (s *server) BidirectionalStreamingMovie(stream pb.Movie_BidirectionalStreamingMovieServer) error {
	for {
		in, err := stream.Recv()
		if err != nil {
			fmt.Printf("server: error receiving from stream: %v\n", err)
			if err == io.EOF {
				return nil
			}
			return err
		}
		fmt.Printf("echoing message %q\n", in.Movie)
		stream.Send(&pb.MovieResponse{Movie: in.Movie})
	}
}
func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *ports))
	if err != nil {
		fmt.Printf("Server not found! %v", err)
	}
	fmt.Printf("server listening at port %v\n", lis.Addr())
	s := grpc.NewServer()
	pb.RegisterMovieServer(s, &server{})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
