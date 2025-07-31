package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	
	pb "github.com/reaksa-maii/one_digital_grpc_getway/proto/echo/v3"
)

var (
	errMassaeg       = status.Errorf(codes.InvalidArgument, "missing invalide")
	errInvalideToken = status.Error(codes.Unauthenticated, "unautheticate token")
)
type server struct {
	
}
func logging(format string, a ...any) {
	fmt.Printf("Log Error: \t"+format+"\n", a...)
}
func validate(authorization []string) bool {
	if len(authorization) < 1 {
		return false
	}
	tokens := strings.TrimPrefix(authorization[0], "Bearer")
	return tokens == "some-secret-token"
}
func unaryIntercepter(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errMassaeg
	}
	if !validate(md["authorization"]) {
		return nil, errInvalideToken
	}
	m, err := handler(ctx, req)
	if err != nil {
		logging("Log gRPC server %v", err)

	}
	return m, err
}

type wrappedStream struct {
	grpc.ServerStream
}

func (w *wrappedStream) RecvMsg(m any) error {
	logging("Receive a message (Type: %T) at %s", m, time.Now().Format(time.RFC3339))
	return w.ServerStream.RecvMsg(m)
}

func (w *wrappedStream) SendMsg(m any) error {
	logging("Send a message (Type: %T) at %v", m, time.Now().Format(time.RFC3339))
	return w.ServerStream.SendMsg(m)
}

func newWrappedStream(s grpc.ServerStream) grpc.ServerStream {
	return &wrappedStream{s}
}
