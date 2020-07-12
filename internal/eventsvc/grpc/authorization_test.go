package grpc

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	authv1pb "github.com/shanvl/garbage/api/auth/v1/pb"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	testInvalidToken = "invalid token"
	testTimeout      = "timeout"
	testUnauthorized = "unauthorized"
	testUnknown      = "testUnknown"
	testUserID       = "some user id"
)

var testUnknownError = status.Error(codes.Unknown, "unknown error")

func TestAuthorization_Authorize(t *testing.T) {
	token := "some token"
	tests := []struct {
		name   string
		method string
		token  string
		userID string
		err    error
	}{
		{
			name:   "invalid token",
			method: testInvalidToken,
			token:  token,
			userID: "",
			err:    ErrInvalidAccessToken,
		},
		{
			name:   "unauthorized",
			method: testUnauthorized,
			token:  token,
			userID: "",
			err:    ErrUnauthorized,
		},
		{
			name:   "unknown error",
			method: testUnknown,
			token:  token,
			userID: "",
			err:    testUnknownError,
		},
		{
			name:   "timeout",
			method: testTimeout,
			token:  token,
			userID: "",
			err:    status.Error(codes.DeadlineExceeded, "context deadline exceeded"),
		},
		{
			name:   "ok",
			method: "some allowed method",
			token:  token,
			userID: testUserID,
			err:    nil,
		},
	}
	srvAddr := startTestAuthServer(t)
	authClient := newTestAuthClient(t, srvAddr)
	authSvc := NewAuthService(authClient, 100*time.Millisecond)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authClaims, err := authSvc.Authorize(context.Background(), tt.token, tt.method)
			if err != nil && !errors.Is(err, tt.err) {
				t.Errorf("Authorize() want err: %v, got: %v", tt.err, err)
			}
			if err != nil && tt.userID != "" {
				t.Errorf("Authorize() err == %v, testUserID: %v", err, authClaims.UserID)
			}
			if err == nil && tt.userID != authClaims.UserID {
				t.Errorf("Authorize() testUserID == %v, want: %v", authClaims.UserID, tt.userID)
			}
		})
	}
}

func newTestAuthClient(t *testing.T, srvAddress string) authv1pb.AuthServiceClient {
	cc, err := grpc.Dial(srvAddress, grpc.WithInsecure())
	require.NoError(t, err)
	if err != nil {
		t.Fatalf("couldn't connect to test auth srv: %v", err)
	}
	return authv1pb.NewAuthServiceClient(cc)
}

func startTestAuthServer(t *testing.T) string {
	t.Helper()
	authServer := grpc.NewServer()
	authv1pb.RegisterAuthServiceServer(authServer, newTestAuthServer())

	l, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("wasn't able to start test auth server: %v", err)
	}
	go authServer.Serve(l)

	return l.Addr().String()
}

func newTestAuthServer() authv1pb.AuthServiceServer {
	return testAuthSvc{}
}

type testAuthSvc struct {
}

func (t testAuthSvc) Authorize(_ context.Context, req *authv1pb.AuthorizeRequest) (*authv1pb.AuthorizeResponse,
	error) {
	if req.GetMethod() == testInvalidToken {
		return nil, status.Error(codes.InvalidArgument, "invalid token")
	}
	if req.GetMethod() == testUnauthorized {
		return nil, status.Error(codes.PermissionDenied, "permission denied")
	}
	if req.GetMethod() == testUnknown {
		return nil, testUnknownError
	}
	if req.GetMethod() == testTimeout {
		time.Sleep(500 * time.Millisecond)
		return nil, testUnknownError
	}
	return &authv1pb.AuthorizeResponse{UserId: testUserID}, nil
}

func (t testAuthSvc) ActivateUser(_ context.Context, _ *authv1pb.ActivateUserRequest) (*empty.Empty, error) {
	return nil, nil
}

func (t testAuthSvc) ChangeUserRole(_ context.Context, _ *authv1pb.ChangeUserRoleRequest) (*empty.Empty, error) {
	return nil, nil
}

func (t testAuthSvc) CreateUser(_ context.Context, _ *authv1pb.CreateUserRequest) (*authv1pb.
	CreateUserResponse, error) {
	return nil, nil
}

func (t testAuthSvc) DeleteUser(_ context.Context, _ *authv1pb.DeleteUserRequest) (*empty.Empty, error) {
	return nil, nil
}

func (t testAuthSvc) FindUser(_ context.Context, _ *authv1pb.FindUserRequest) (*authv1pb.FindUserResponse,
	error) {
	return nil, nil
}

func (t testAuthSvc) FindUsers(_ context.Context, _ *authv1pb.FindUsersRequest) (*authv1pb.FindUsersResponse,
	error) {
	return nil, nil
}

func (t testAuthSvc) Login(_ context.Context, _ *authv1pb.LoginRequest) (*authv1pb.LoginResponse, error) {
	return nil, nil
}

func (t testAuthSvc) Logout(_ context.Context, _ *authv1pb.LogoutRequest) (*empty.Empty, error) {
	return nil, nil
}

func (t testAuthSvc) LogoutAllClients(_ context.Context, _ *empty.Empty) (*empty.Empty, error) {
	return nil, nil
}

func (t testAuthSvc) RefreshTokens(_ context.Context, _ *authv1pb.RefreshTokensRequest) (*authv1pb.
	RefreshTokensResponse, error) {
	return nil, nil
}

func newTestAuthService() AuthorizationService {
	return &testAuthService{}
}

type testAuthService struct{}

func (t testAuthService) Authorize(_ context.Context, _, _ string) (*AuthClaims, error) {
	return &AuthClaims{
		ClientID: "clientid",
		UserID:   "userid",
	}, nil
}
