package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/reaksa-maii/one_digital_grpc_getway/proto/book/v3"
)

var port = flag.Int("port", 50051, "server port")

type server struct {
	pb.UnimplementedBookServer
}

func (s *server) UnaryBook(_ context.Context, in *pb.BookRequest) (*pb.BookResponse, error) {
	fmt.Printf("Unary Book Massage: %q\n", in.Book)
	return &pb.BookResponse{Book: in.Book}, nil
}

func main() {

	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		fmt.Printf("server notfound from: %v", err)
	}

	fmt.Printf("server listening at port %v\n", lis.Addr())

	s := grpc.NewServer()
	pb.RegisterBookServer(s, &server{})
	reflection.Register(s)
	
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed server start: %v", err)
	}
}
