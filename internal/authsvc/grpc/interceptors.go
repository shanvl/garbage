package grpc

import (
	"context"
	"errors"

	"github.com/shanvl/garbage/internal/authsvc"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// grpc_recovery interceptor helper function. Handle panic by logging it and returning Internal error
func (s *Server) handleRecovery(p interface{}) error {
	s.log.Error("panic triggered", zap.Any("panic", p))
	return status.Error(codes.Internal, "internal server error")
}

const AuthCtxKey = "auth"

func authClaimsFromCtx(ctx context.Context) (*authsvc.UserClaims, error) {
	creds, ok := ctx.Value(AuthCtxKey).(*authsvc.UserClaims)
	if !ok {
		return nil, errors.New("no auth claims in ctx")
	}
	return creds, nil
}

func authClaimsToCtx(ctx context.Context, claims *authsvc.UserClaims) context.Context {
	return context.WithValue(ctx, AuthCtxKey, claims)
}
