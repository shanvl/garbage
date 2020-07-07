package grpc

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	authv1pb "github.com/shanvl/garbage/api/auth/v1/pb"
)

// repo.GetUserByEmail
// repo.StoreClient
// repo.DeleteClient
// repo.DeleteUserClients
func (s *Server) Authorize(ctx context.Context, request *authv1pb.AuthorizeRequest) (*authv1pb.AuthorizeResponse, error) {
	// parse the token
	// get the user from the token
	// check the role and decide whether the user can access rpc

	// call s.authorize where all the parsing and decision making happens
	panic("implement me")
}

func (s *Server) Login(ctx context.Context, request *authv1pb.LoginRequest) (*authv1pb.LoginResponse, error) {
	// get the user
	// check if the user is active
	// check if password == u.password_hash
	// generate new tokens and client id
	// save the tokens
	// send the tokens to the user

	// repo.GetUserByID
	// repo.Login OR repo.StoreClient
	panic("implement me")
}

func (s *Server) Logout(ctx context.Context, request *authv1pb.LogoutRequest) (*empty.Empty, error) {
	// get clientID from the token
	// remove it from db along with refresh toke

	// repo.Logout OR repo.DeleteClient
	panic("implement me")
}

func (s *Server) LogoutAllClients(ctx context.Context, empty *empty.Empty) (*empty.Empty, error) {
	// get userID from the token
	// delete all clients of that user

	// repo.LogoutAllUserClients OR repo.DeleteUserClients
	panic("implement me")
}
