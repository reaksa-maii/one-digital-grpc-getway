package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"io"

	pb "github.com/reaksa-maii/one_digital_grpc_getway/proto/movie/v3"
	// "github.com/reaksa-maii/one_digital_grpc_getway/utilities"
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

func (s *server) BidirectionalStreamingMovie(stream pb.Movie_BidirectionalStreamingMovieServer) error {
	for {
		in, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			fmt.Printf("server: error receiving from stream: %v\n", err)
			return err
		}
		fmt.Printf("bidi echoing message %q\n", in.MovieSize)
		stream.Send(&pb.MovieResponse{MovieSize: in.MovieSize})
	}
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterMovieServer(s, &server{})
	reflection.Register(s)

	fmt.Printf("server listening at %v\n", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
