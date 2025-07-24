package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/reaksa-maii/one_digital_grpc_getway/proto/echo/v1"
)

var port = flag.Int("port", 50051, "the port to serve on")

type server struct {
	pb.UnimplementedEchoServer
}

func (s *server) BidirectionalStreamingEcho(stream pb.Echo_BidirectionalStreamingEchoServer) error {
	for {
		in, err := stream.Recv()
		if err != nil {
			fmt.Printf("server: error receiving from stream: %v\n", err)
			if err == io.EOF {
				return nil
			}
			return err
		}
		fmt.Printf("echoing message %q\n", in.Message)
		stream.Send(&pb.EchoResponse{Message: in.Message})
	}
}

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Printf("server listening at port %v\n", lis.Addr())

	s := grpc.NewServer()

	pb.RegisterEchoServer(s, &server{})

	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
