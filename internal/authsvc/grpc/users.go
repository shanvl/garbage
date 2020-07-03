package grpc

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	authv1pb "github.com/shanvl/garbage/api/auth/v1/pb"
)

// repo.GetUserByID
// repo.GetUsers
// repo.StoreUser
// repo.DeleteUserByID
func (s *Server) ActivateUser(ctx context.Context, request *authv1pb.ActivateUserRequest) (*empty.Empty, error) {
	// get the user
	// compare the activate_tokens
	// change its active field to true
	// store the user

	// repo.GetUserByID
	// repo.StoreUser
	panic("implement me")
}

func (s *Server) ChangeUserRole(ctx context.Context, request *authv1pb.ChangeUserRoleRequest) (*empty.Empty, error) {
	// get the user
	// change its role
	// store the user

	// repo.StoreUser
	panic("implement me")
}

func (s *Server) CreateUser(ctx context.Context, request *authv1pb.CreateUserRequest) (*authv1pb.CreateUserResponse, error) {
	// generate activate_token
	// generate userID
	// store the user

	// repo.StoreUser
	panic("implement me")
}

func (s *Server) DeleteUser(ctx context.Context, request *authv1pb.DeleteUserRequest) (*empty.Empty, error) {
	// delete the user

	// repo.DeleteUserByID
	panic("implement me")
}

func (s *Server) FindUser(ctx context.Context, request *authv1pb.FindUserRequest) (*authv1pb.FindUserResponse, error) {
	// find user by id

	// repo.GetUserByID
	panic("implement me")
}

func (s *Server) FindUsers(ctx context.Context, request *authv1pb.FindUsersRequest) (*authv1pb.FindUsersResponse, error) {
	// find users via text_search

	// repo.GetUsers
	panic("implement me")
}
