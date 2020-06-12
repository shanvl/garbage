package grpc

import (
	"log"

	grpcRecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Interceptor interface {
	Unary() grpc.UnaryServerInterceptor
	Stream() grpc.StreamServerInterceptor
}

type PanicInterceptor func(p interface{}) (err error)

func NewPanicInterceptor() PanicInterceptor {
	return func(p interface{}) (err error) {
		log.Printf("panic triggered: %v", p)
		return status.Error(codes.Internal, "internal server error")
	}
}

func (i PanicInterceptor) Unary() grpc.UnaryServerInterceptor {
	return grpcRecovery.UnaryServerInterceptor(grpcRecovery.WithRecoveryHandler(grpcRecovery.RecoveryHandlerFunc(i)))
}

func (i PanicInterceptor) Stream() grpc.StreamServerInterceptor {
	return grpcRecovery.StreamServerInterceptor(grpcRecovery.WithRecoveryHandler(grpcRecovery.RecoveryHandlerFunc(i)))
}
