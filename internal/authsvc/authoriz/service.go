// package authoriz is responsible for authorization of the users' requests
package authoriz

import (
	"context"
	"errors"
	"fmt"

	"github.com/shanvl/garbage/internal/authsvc"
	"github.com/shanvl/garbage/pkg/valid"
)

var ErrUnauthorized = errors.New("unauthorized")

// Service is responsible for authorization of the users' requests
type Service interface {
	// Authorize decides whether the user has access to the requested RPC
	Authorize(ctx context.Context, accessToken, method string) (authsvc.UserClaims, error)
}

type service struct {
	tm           authsvc.TokenManager
	protectedRPC map[string][]authsvc.Role
}

func NewService(tm authsvc.TokenManager, protectedRPC map[string][]authsvc.Role) Service {
	return &service{tm, protectedRPC}
}

// Authorize decides whether the user has access to the requested RPC
func (s *service) Authorize(_ context.Context, accessToken, method string) (authsvc.UserClaims, error) {
	// validate the arguments
	if method == "" {
		return authsvc.UserClaims{}, valid.NewError("method", "method is required")
	}
	// if the method is not protected, it doesn't need an access token. A method isn't protected if it isn't in the map
	if _, ok := s.protectedRPC[method]; !ok {
		return authsvc.UserClaims{}, nil
	}
	// verify the token and extract its claims
	claims, err := s.tm.Verify(accessToken)
	if err != nil {
		return authsvc.UserClaims{}, fmt.Errorf("%w: %v", authsvc.ErrInvalidAccessToken, err)
	}
	// convert string role from the claims to authsvc.Role
	role, err := authsvc.StringToRole(claims.Role)
	if err != nil {
		return authsvc.UserClaims{}, fmt.Errorf("%w: invalid role: %s: %v", authsvc.ErrInvalidAccessToken, role, err)
	}
	// check whether the user's role has access to the method
	roles, ok := s.protectedRPC[method]
	if ok {
		for _, r := range roles {
			if r == role {
				return claims, nil
			}
		}
		return authsvc.UserClaims{}, ErrUnauthorized
	}
	return claims, nil
}

// ProtectedRPCMap creates a map of protected RPCs. Later it can be changed to more elaborate procedure involving db,
// but for now it's a simple map in memory
func ProtectedRPCMap() map[string][]authsvc.Role {
	const authSvcPrefix = "/shanvl.garbage.auth.v1.AuthService/"
	const eventSvcPrefix = "/shanvl.garbage.events.v1.EventsService/"
	return map[string][]authsvc.Role{
		authSvcPrefix + "ChangeUserRole":        {authsvc.Admin, authsvc.Root},
		authSvcPrefix + "CreateUser":            {authsvc.Admin, authsvc.Root},
		authSvcPrefix + "DeleteUser":            {authsvc.Admin, authsvc.Root},
		authSvcPrefix + "FindUser":              {authsvc.Admin, authsvc.Member, authsvc.Root},
		authSvcPrefix + "FindUsers":             {authsvc.Admin, authsvc.Member, authsvc.Root},
		authSvcPrefix + "Logout":                {authsvc.Admin, authsvc.Member, authsvc.Root},
		authSvcPrefix + "LogoutAllClients":      {authsvc.Admin, authsvc.Member, authsvc.Root},
		authSvcPrefix + "RefreshTokens":         {authsvc.Admin, authsvc.Member, authsvc.Root},
		eventSvcPrefix + "AddPupils":            {authsvc.Admin, authsvc.Root},
		eventSvcPrefix + "ChangePupilClass":     {authsvc.Admin, authsvc.Root},
		eventSvcPrefix + "ChangePupilResources": {authsvc.Admin, authsvc.Member, authsvc.Root},
		eventSvcPrefix + "CreateEvent":          {authsvc.Admin, authsvc.Root},
		eventSvcPrefix + "DeleteEvent":          {authsvc.Admin, authsvc.Root},
		eventSvcPrefix + "FindClasses":          {authsvc.Admin, authsvc.Member, authsvc.Root},
		eventSvcPrefix + "FindEvents":           {authsvc.Admin, authsvc.Member, authsvc.Root},
		eventSvcPrefix + "FindEventByID":        {authsvc.Admin, authsvc.Member, authsvc.Root},
		eventSvcPrefix + "FindEventClasses":     {authsvc.Admin, authsvc.Member, authsvc.Root},
		eventSvcPrefix + "FindEventPupils":      {authsvc.Admin, authsvc.Member, authsvc.Root},
		eventSvcPrefix + "FindEventPupilByID":   {authsvc.Admin, authsvc.Member, authsvc.Root},
		eventSvcPrefix + "FindPupils":           {authsvc.Admin, authsvc.Member, authsvc.Root},
		eventSvcPrefix + "RemovePupils":         {authsvc.Admin, authsvc.Root},
	}
}
