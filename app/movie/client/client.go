package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	pb "github.com/reaksa-maii/one_digital_grpc_getway/proto/movie/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

var addr = flag.String("addr", "localhost:50051", "the address to connect to")

func sendMassage(stream pb.Movie_BidirectionalStreamingMovieClient, msg string) error {
	fmt.Printf("sending message %q\n", msg)
	return stream.Send(&pb.MovieRequest{Movie: msg})
}
func recvMassage(stream pb.Movie_BidirectionalStreamingMovieClient, ErrorCode codes.Code) {
	res, err := stream.Recv()
	if status.Code(err) != ErrorCode {
		log.Fatalf("stream.Recv() = %v, %v; want _, status.Code(err)=%v", res, err, ErrorCode)
	}
	if err != nil {
		fmt.Printf("stream.Recv() returned expected error %v\n", err)
		return
	}
	fmt.Printf("receving massage %q\n", res.GetMovie())
}
func main() {
	flag.Parse()
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not to connect server %v", err)
	}
	defer conn.Close()
	c := pb.NewMovieClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	stream, err := c.BidirectionalStreamingMovie(ctx)
	if err != nil {
		log.Fatalf("error stremig data from %v", err)
	}
	if err := sendMassage(stream, "Hi"); err != nil {
		log.Fatalf("Error sending adta %v", err)
	}
	recvMassage(stream, codes.OK)
	recvMassage(stream, codes.OK)
	fmt.Printf("stream are context")
	cancel()
	sendMassage(stream, "closed")
	recvMassage(stream, codes.Canceled)
}
