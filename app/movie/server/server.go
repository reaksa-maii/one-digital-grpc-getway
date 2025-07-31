package server

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/reaksa-maii/one_digital_grpc_getway/proto/movie/v1"
)

var (
	ports = flag.Int("port", 50052, "port server run time")
	errMassaeg = status.Errorf(codes.InvalidArgument, "missing matada")
	errInvalideToken = status.Errorf(codes.Unauthenticated, "invalide token")
)