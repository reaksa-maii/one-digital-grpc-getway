package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/reaksa-maii/one_digital_grpc_getway/proto/book/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var port = flag.Int("port", 50051, "the port to serve on")

type server struct {
	pb.UnimplementedBookServer
}

func (s *server) UnaryBook(ctx context.Context, in *pb.BookRequest) (*pb.BookResponse, error) {
	fmt.Printf("Book Unary Service")
	return &pb.BookResponse{Book: in.Book}, nil
}
func main() {

	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Printf("server listening at %v\n", lis.Addr())

	s := grpc.NewServer()
	pb.RegisterBookServer(s, &server{})
	reflection.Register(s)
	s.Serve(lis)
}
