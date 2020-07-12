package grpc

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	authv1pb "github.com/shanvl/garbage/api/auth/v1/pb"
)

// Login generates, saves and returns auth credentials for the user if the given password and the email are correct
func (s *Server) Login(ctx context.Context, req *authv1pb.LoginRequest) (*authv1pb.LoginResponse, error) {
	user, creds, err := s.authentSvc.Login(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, s.handleError(err)
	}

	return &authv1pb.LoginResponse{User: userToProto(user), Tokens: credsToProto(creds)}, nil
}

// Logout deletes the user's client and refresh token from the db, thus, logging the user out
func (s *Server) Logout(ctx context.Context, _ *authv1pb.LogoutRequest) (*empty.Empty, error) {
	// claims are put to ctx by auth interceptor
	claims, err := authClaimsFromCtx(ctx)
	if err != nil {
		return nil, s.handleError(err)
	}
	err = s.authentSvc.Logout(ctx, claims.ClientID)
	if err != nil {
		return nil, s.handleError(err)
	}
	return &empty.Empty{}, nil
}

// LogoutAllClients deletes all the user's clients and refresh tokens, thus, logging the user out of every device
func (s *Server) LogoutAllClients(ctx context.Context, _ *empty.Empty) (*empty.Empty, error) {
	// claims are put to ctx by auth interceptor
	claims, err := authClaimsFromCtx(ctx)
	if err != nil {
		return nil, s.handleError(err)
	}
	err = s.authentSvc.LogoutAllClients(ctx, claims.Subject)
	if err != nil {
		return nil, s.handleError(err)
	}
	return &empty.Empty{}, nil
}

// RefreshTokens verifies the given refresh token and then creates, saves and returns new auth credentials
func (s *Server) RefreshTokens(ctx context.Context, req *authv1pb.RefreshTokensRequest) (*authv1pb.
	RefreshTokensResponse, error) {

	creds, err := s.authentSvc.RefreshTokens(ctx, req.GetRefreshToken())
	if err != nil {
		return nil, s.handleError(err)
	}

	return &authv1pb.RefreshTokensResponse{Tokens: credsToProto(creds)}, nil
}
