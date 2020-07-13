package grpc

import (
	"context"
	"errors"
	"fmt"
	"time"

	authv1pb "github.com/shanvl/garbage/api/auth/v1/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrInvalidAccessToken = errors.New("invalid token")
var ErrUnauthorized = errors.New("no permission to access this RPC")

// AuthorizationService is used to determine whether the user has access to the requested RPC
type AuthorizationService interface {
	Authorize(ctx context.Context, token, method string) (*AuthClaims, error)
}

// authService is an implementation of AuthorizationService which uses a gRPC client to call separate auth server
type authService struct {
	svc authv1pb.AuthServiceClient
	// time to wait for the auth svc response
	timeout time.Duration
}

// NewAuthService returns authService
func NewAuthService(pbAuthSvc authv1pb.AuthServiceClient, timeout time.Duration) AuthorizationService {
	return &authService{pbAuthSvc, timeout}
}

// Authorize requests the permission to use one of the eventsvc methods and returns the user's id if it gets the
// permission
func (a *authService) Authorize(ctx context.Context, token, method string) (*AuthClaims, error) {
	ctxWithDeadline, cancel := context.WithDeadline(ctx, time.Now().Add(a.timeout))
	defer cancel()

	resp, err := a.svc.Authorize(ctxWithDeadline, &authv1pb.AuthorizeRequest{
		Method: method,
		Token:  token,
	})
	// convert gRPC specific errors to domain errors
	if err != nil {
		grpcErr := status.Convert(err)
		switch grpcErr.Code() {
		case codes.PermissionDenied:
			return nil, fmt.Errorf("%w: %v", ErrUnauthorized, grpcErr.Message())
		case codes.InvalidArgument:
			fallthrough
		case codes.Unauthenticated:
			return nil, fmt.Errorf("%w: %v", ErrInvalidAccessToken, grpcErr.Message())
		default:
			return nil, err
		}
	}
	return &AuthClaims{UserID: resp.GetUserId(), ClientID: resp.GetClientId()}, nil
}

// AuthClaims contain some info about the request
type AuthClaims struct {
	ClientID string
	UserID   string
}
