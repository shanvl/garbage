package grpc

import (
	authv1pb "github.com/shanvl/garbage/api/auth/v1/pb"
	"github.com/shanvl/garbage/internal/authsvc"
	"github.com/shanvl/garbage/internal/authsvc/authent"
	"github.com/shanvl/garbage/internal/authsvc/users"
)

var roleProtoMap = map[authsvc.Role]authv1pb.Role{
	authsvc.Admin:  authv1pb.Role_ROLE_ADMIN,
	authsvc.Member: authv1pb.Role_ROLE_MEMBER,
	authsvc.Root:   authv1pb.Role_ROLE_ROOT,
}

var protoRoleMap = map[authv1pb.Role]authsvc.Role{
	authv1pb.Role_ROLE_ADMIN:  authsvc.Admin,
	authv1pb.Role_ROLE_MEMBER: authsvc.Member,
	authv1pb.Role_ROLE_ROOT:   authsvc.Root,
}

var protoUserSortingMap = map[authv1pb.UserSorting]users.Sorting{
	authv1pb.UserSorting_USER_SORTING_NAME_ASC:   users.NameAsc,
	authv1pb.UserSorting_USER_SORTING_NAME_DESC:  users.NameDes,
	authv1pb.UserSorting_USER_SORTING_EMAIL_ASC:  users.EmailAsc,
	authv1pb.UserSorting_USER_SORTING_EMAIL_DESC: users.EmailDes,
	authv1pb.UserSorting_USER_SORTING_UNKNOWN:    users.Unspecified,
}

// protoToRole converts authv1pb.Role to authsvc.Role
func protoToRole(proto authv1pb.Role) (authsvc.Role, error) {
	role, ok := protoRoleMap[proto]
	if !ok {
		return 0, authsvc.ErrUnknownRole
	}
	return role, nil
}

// userToProto converts *authsvc.User to *authv1pb.User
func userToProto(user *authsvc.User) *authv1pb.User {
	return &authv1pb.User{
		Id:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      roleProtoMap[user.Role],
	}
}

// credsToProto converts authent.Creds to *authv1pb.Tokens
func credsToProto(creds authent.AuthCreds) *authv1pb.Tokens {
	return &authv1pb.Tokens{
		AccessToken:  creds.Access,
		RefreshToken: creds.Refresh,
		ClientId:     creds.ClientID,
	}
}
