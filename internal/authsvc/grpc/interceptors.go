package grpc

import (
	"context"
	"errors"
	"strings"

	"github.com/shanvl/garbage/internal/authsvc"
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

		// get access token from auth header
		token, err := getAccessTokenFromAuthHeader(ctx, "bearer")
		if err != nil {
			return nil, s.handleError(err)
		}

		// call the authorization service
		claims, err := s.authoriz.Authorize(ctx, token, info.FullMethod)
		if err != nil {
			return nil, s.handleError(err)
		}

		// add the claims to the ctx
		ctx = authClaimsToCtx(ctx, claims)

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

		// get access token from auth header
		token, err := getAccessTokenFromAuthHeader(stream.Context(), "bearer")
		if err != nil {
			return s.handleError(err)
		}

		// call the authorization service
		claims, err := s.authoriz.Authorize(stream.Context(), token, info.FullMethod)
		if err != nil {
			return s.handleError(err)
		}

		// add the claims to the ctx
		stream = newStreamWithAuthCtx(claims, stream)

		return handler(srv, stream)
	}
}

// ServerStream wrapper used in adding auth claims to ctx
type streamWithAuthCtx struct {
	claims *authsvc.UserClaims
	grpc.ServerStream
}

// Context populates ctx with auth claims
func (s *streamWithAuthCtx) Context() context.Context {
	return authClaimsToCtx(s.ServerStream.Context(), s.claims)
}

func newStreamWithAuthCtx(claims *authsvc.UserClaims, s grpc.ServerStream) grpc.ServerStream {
	return &streamWithAuthCtx{
		claims:       claims,
		ServerStream: s,
	}
}

// getAccessTokenFromAuthHeader extracts an access token from the ctx
func getAccessTokenFromAuthHeader(ctx context.Context, scheme string) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", authsvc.ErrInvalidAccessToken
	}
	authHeader := md["authorization"]
	if len(authHeader) < 1 {
		return "", authsvc.ErrInvalidAccessToken
	}
	splits := strings.SplitN(authHeader[0], " ", 2)
	if len(splits) < 2 {
		return "", authsvc.ErrInvalidAccessToken
	}
	if !strings.EqualFold(splits[0], scheme) {
		return "", authsvc.ErrInvalidAccessToken
	}
	return splits[1], nil
}

const AuthCtxKey = "auth"

// authClaimsFromCtx extracts auth claims from ctx
func authClaimsFromCtx(ctx context.Context) (*authsvc.UserClaims, error) {
	creds, ok := ctx.Value(AuthCtxKey).(*authsvc.UserClaims)
	if !ok {
		return nil, errors.New("no auth claims in ctx")
	}
	return creds, nil
}

// authClaimsToCtx adds auth claims to ctx
func authClaimsToCtx(ctx context.Context, claims *authsvc.UserClaims) context.Context {
	return context.WithValue(ctx, AuthCtxKey, claims)
}
