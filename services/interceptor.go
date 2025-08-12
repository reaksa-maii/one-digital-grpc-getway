package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"

	mb "github.com/reaksa-maii/one_digital_grpc_getway/proto/movie/v3"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

var (
	port = flag.Int("port", 50051, "the port to serve on")

	errMetadata = status.Errorf(codes.InvalidArgument, "missing metadata")
	errToken    = status.Errorf(codes.Unauthenticated, "invalid token")
)

func logger(format string, a ...any) {
	fmt.Printf("LOG:\t"+format+"\n", a...)
}

type InterceptServer struct {
	mb.UnimplementedMovieServer
}

func (s *InterceptServer) UnaryMovie(ctx context.Context, req *mb.MovieRequest) (*mb.MovieResponse, error) {
	fmt.Printf("unary movie message %q\n", req.Title)
	return &mb.MovieResponse{Title: req.Title}, nil
}
func (s *InterceptServer) BidirectionalStreamingMovie(stream mb.Movie_BidirectionalStreamingMovieServer) error {
	for {
		in, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			fmt.Printf("server: error receiving from stream: %v\n", err)
			return err
		}
		fmt.Printf("bidi echoing message %q\n", in.Title)
		stream.Send(&mb.MovieResponse{Title: in.Title})
	}
}

func valid(authorization []string) bool {
	if len(authorization) < 1 {
		return false
	}
	token := strings.TrimPrefix(authorization[0], "Bearer ")
	return token == "some-secret-token"
}
func unaryInterceptor(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	// authentication (token verification)
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errMetadata
	}
	if !valid(md["authorization"]) {
		return nil, errToken
	}
	m, err := handler(ctx, req)
	if err != nil {
		logger("RPC failed with error: %v", err)
	}
	return m, err
}

type wrappedStreaming struct {
	grpc.ServerStream
}

func (w *wrappedStreaming) RecvMsg(m any) error {
	logger("Receive a message (Type: %T) at %s", m, time.Now().Format(time.RFC3339))
	return w.ServerStream.RecvMsg(m)
}

func (w *wrappedStreaming) SendMsg(m any) error {
	logger("Send a message (Type: %T) at %v", m, time.Now().Format(time.RFC3339))
	return w.ServerStream.SendMsg(m)
}

func newWrappedStream(s grpc.ServerStream) grpc.ServerStream {
	return &wrappedStreaming{s}
}
func streamInterceptor(srv any, ss grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	md, ok := metadata.FromIncomingContext(ss.Context())
	if !ok {
		return errMetadata
	}
	if !valid(md["authorization"]) {
		return errToken
	}

	err := handler(srv, newWrappedStream(ss))
	if err != nil {
		logger("RPC failed with error: %v", err)
	}
	return err
}

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	creds, err := credentials.NewServerTLSFromFile(
		"../certificate/server-cert.pem",
		"../certificate/server-key.pem",
	)
	if err != nil {
		log.Fatalf("failed to create credentials: %v", err)
	}

	s := grpc.NewServer(grpc.Creds(creds), grpc.UnaryInterceptor(unaryInterceptor), grpc.StreamInterceptor(streamInterceptor))
	// Register EchoServer on the server.
	reflection.Register(s)
	mb.RegisterMovieServer(s, &InterceptServer{})

	fmt.Printf("gRPC server listening at %v\n", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
