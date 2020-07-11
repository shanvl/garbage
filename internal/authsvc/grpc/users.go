package grpc

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	authv1pb "github.com/shanvl/garbage/api/auth/v1/pb"
)

// ActivateUser changes the active state of the user to active and populates it with the provided additional info
func (s *Server) ActivateUser(ctx context.Context, req *authv1pb.ActivateUserRequest) (*empty.Empty, error) {
	_, err := s.users.ActivateUser(ctx, req.GetActivationToken(), req.GetFirstName(), req.GetLastName(),
		req.GetPassword())

	if err != nil {
		return nil, s.handleError(err)
	}
	return &empty.Empty{}, nil
}

// ChangeUserRole changes the user's role to the provided role
func (s *Server) ChangeUserRole(ctx context.Context, req *authv1pb.ChangeUserRoleRequest) (*empty.Empty, error) {
	role, err := protoToRole(req.GetRole())
	if err != nil {
		return nil, s.handleError(err)
	}
	err = s.users.ChangeUserRole(ctx, req.GetId(), role)
	if err != nil {
		return nil, s.handleError(err)
	}
	return &empty.Empty{}, nil
}

// CreateUser creates and stores a user, which must then be activated with the returned activation token
// Note, that the user's password is not needed here, it is required on the activation step
func (s *Server) CreateUser(ctx context.Context, req *authv1pb.CreateUserRequest) (*authv1pb.CreateUserResponse,
	error) {

	userID, activationToken, err := s.users.CreateUser(ctx, req.GetEmail())
	if err != nil {
		return nil, s.handleError(err)
	}

	return &authv1pb.CreateUserResponse{Id: userID, ActivationToken: activationToken}, nil
}

// DeleteUser deletes the user
func (s *Server) DeleteUser(ctx context.Context, req *authv1pb.DeleteUserRequest) (*empty.Empty, error) {
	err := s.users.DeleteUser(ctx, req.GetId())
	if err != nil {
		return nil, s.handleError(err)
	}
	return &empty.Empty{}, nil
}

// FindUser returns the user with the specified id
func (s *Server) FindUser(ctx context.Context, req *authv1pb.FindUserRequest) (*authv1pb.FindUserResponse, error) {
	user, err := s.users.UserByID(ctx, req.GetId())
	if err != nil {
		return nil, s.handleError(err)
	}
	return &authv1pb.FindUserResponse{User: userToProto(user)}, nil
}

// FindUsers returns a sorted list of users
// "nameAndEmail" may consist of any combination of the email, first name and last name parts
func (s *Server) FindUsers(ctx context.Context, req *authv1pb.FindUsersRequest) (*authv1pb.FindUsersResponse, error) {
	users, total, err := s.users.Users(
		ctx,
		req.GetNameAndEmail(),
		protoUserSortingMap[req.GetSorting()],
		int(req.GetAmount()),
		int(req.GetSkip()),
	)
	if err != nil {
		return nil, s.handleError(err)
	}
	// convert []*authsvc.User to []*authv1pb.User
	usersProto := make([]*authv1pb.User, len(users))
	for i, user := range users {
		usersProto[i] = userToProto(user)
	}
	return &authv1pb.FindUsersResponse{Users: usersProto, Total: uint32(total)}, nil
}
