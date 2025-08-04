package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"

	pb "github.com/reaksa-maii/one_digital_grpc_getway/proto/movie/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var ports = flag.Int("port", 50051, "port server run time")

type server struct {
	pb.UnimplementedMovieServer
}

func (s *server) UnaryMovie(_ context.Context, in *pb.MovieRequest) (*pb.MovieResponse, error) {
	fmt.Printf("Unary Movie Massage: %q\n", in.Title)
	return &pb.MovieResponse{Title: in.Title}, nil
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
		fmt.Printf("echoing message %q\n", in.Title)
		stream.Send(&pb.MovieResponse{Title: in.Title})
	}
}
func runGRPCServer() error {
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
func runRESTServer() error {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	pb.RegisterMovieServer(s, &server{})

	// Enable reflection to allow clients to query the server's services
	reflection.Register(s)

	fmt.Println("Starting gRPC server on :8080...")
	if err := s.Serve(lis); err != nil {
		return err
	}

	return nil
}
func main() {
	go func() {
		if err := runRESTServer(); err != nil {
			log.Fatalf("Failed to run REST server: %v", err)
		}
	}()

	if err := runGRPCServer(); err != nil {
		log.Fatalf("Failed to run gRPC server: %v", err)
	}
}
