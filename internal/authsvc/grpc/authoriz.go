package grpc

import (
	"context"

	authv1pb "github.com/shanvl/garbage/api/auth/v1/pb"
)

// Authorize decides whether the user has access to the requested RPC
func (s *Server) Authorize(ctx context.Context, req *authv1pb.AuthorizeRequest) (*authv1pb.AuthorizeResponse, error) {
	details, err := s.authoriz.Authorize(req.GetToken(), req.GetMethod())
	if err != nil {
		return nil, s.handleError(err)
	}
	return &authv1pb.AuthorizeResponse{UserId: details.Subject}, nil
}