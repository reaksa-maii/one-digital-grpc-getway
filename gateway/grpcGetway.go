package gateway

import (
	"fmt"
	"net"

	pb "github.com/reaksa-maii/one_digital_grpc_getway/proto/movie/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	pb.UnimplementedMovieServer
}

func RungRPC() error {
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	pb.RegisterMovieServer(s, &server{})

	// Enable reflection to allow clients to query the server's services
	reflection.Register(s)

	fmt.Println("Starting gRPC server on localhost:50051...")
	if err := s.Serve(lis); err != nil {
		return err
	}

	return nil
}
