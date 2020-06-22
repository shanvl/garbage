package grpc

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// grpc_zap interceptor helper function. When to log payload: always
func (s *Server) payloadDecider(context.Context, string, interface{}) bool {
	return true
}

// grpc_recovery interceptor helper function. Handle panic by logging it and returning Internal error
func (s *Server) handleRecovery(p interface{}) error {
	s.log.Error("panic triggered", zap.Any("panic", p))
	return status.Error(codes.Internal, "internal server error")
}
