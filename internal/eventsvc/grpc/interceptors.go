package grpc

import (
	"context"
	"errors"
	"strings"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// grpc_recovery interceptor helper function. Handle panic by logging it and returning Internal error
func (s *Server) handleRecovery(p interface{}) error {
	s.log.Error("panic triggered", zap.Any("panic", p))
	return status.Error(codes.Internal, "internal server error")
}

// authUnaryInterceptor talks to the auth service to get the permission to access the rpc
func (s *Server) authUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		err := s.authorize(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

// authServerInterceptor talks to the auth service to get the permission to access the rpc
func (s *Server) authStreamInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {

		err := s.authorize(stream.Context(), info.FullMethod)
		if err != nil {
			return err
		}

		return handler(srv, stream)
	}
}

// authorize talks to the auth service to get the permission to access the rpc
func (s *Server) authorize(ctx context.Context, method string) error {
	// get auth header
	token, err := getAuthTokenFromCtx(ctx, "bearer")
	if err != nil {
		return status.Error(codes.Unauthenticated, err.Error())
	}

	// call the service
	userID, err := s.authSvc.Authorize(ctx, method, token)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidToken):
			return status.Errorf(codes.Unauthenticated, err.Error())
		case errors.Is(err, ErrPermissionDenied):
			return status.Errorf(codes.PermissionDenied, err.Error())
		default:
			return s.handleError(err)
		}
	}
	s.log.Info(userID)

	return nil
}

// getAuthTokenFromCtx extracts an access token from the ctx
func getAuthTokenFromCtx(ctx context.Context, scheme string) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", ErrInvalidToken
	}
	authHeader := md["authorization"]
	if len(authHeader) < 1 {
		return "", ErrInvalidToken
	}
	splits := strings.SplitN(authHeader[0], " ", 2)
	if len(splits) < 2 {
		return "", ErrInvalidToken
	}
	if !strings.EqualFold(splits[0], scheme) {
		return "", ErrInvalidToken
	}
	return splits[1], nil
}
