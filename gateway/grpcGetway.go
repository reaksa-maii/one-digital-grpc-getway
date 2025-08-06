package gateway

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/reaksa-maii/one_digital_grpc_getway/proto/podcast/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var port = flag.Int("port", 50051, "the port to serve on")

type server struct {
	pb.UnimplementedPodcatServiceServer
}

func (s *server) CreatePodcast(ctx context.Context, in *pb.PodcatRequest) (*pb.PodcatResponse, error) {
	return &pb.PodcatResponse{
		Id:          in.Id,
		PodcatSize:  in.PodcatSize,
		Title:       in.Title,
		Category:    in.Category,
		Description: in.Description,
		Duration:    in.Duration,
	}, nil
}

func (s *server) GetPodcatByTitle(ctx context.Context, req *pb.PodcatRequest) (*pb.PodcatResponse, error) {
	log.Printf("GetPodcatByTitle called with: %s\n", req.Title)
	return &pb.PodcatResponse{Title: req.GetTitle()}, nil
}

func (s *server) UnaryPodcast(ctx context.Context, req *pb.PodcatRequest) (*pb.PodcatResponse, error) {
	log.Printf("GetPodcatByTitle called with: %s\n", req.Title)
	return &pb.PodcatResponse{Title: req.GetTitle()}, nil
}

func (s *server) ServerStreamingPodcat(req *pb.PodcatRequest, stream pb.PodcatService_ServerStreamingPodcatServer) error {
	for i := 0; i < 5; i++ {
		stream.Send(&pb.PodcatResponse{Title: req.Title, PodcatSize: "size", Duration: float64(i)})
	}
	return nil
}

func (s *server) ClientStreamingPodcat(stream pb.PodcatService_ClientStreamingPodcatServer) error {
	var last *pb.PodcatRequest
	for {
		req, err := stream.Recv()
		if err != nil {
			break
		}
		last = req
	}
	return stream.SendAndClose(&pb.PodcatResponse{Title: last.GetTitle(), Description: "Last received"})
}

func (s *server) BidirectionalStreamingPodcat(stream pb.PodcatService_BidirectionalStreamingPodcatServer) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			break
		}
		resp := &pb.PodcatResponse{Title: req.Title + " (echo)"}
		stream.Send(resp)
	}
	return nil
}
func RungRPC() error {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	pb.RegisterPodcatServiceServer(s, &server{})

	reflection.Register(s)
	fmt.Printf("gRPC server listening at %v\n", lis.Addr())
	if err := s.Serve(lis); err != nil {
		return err
	}

	return nil
}
