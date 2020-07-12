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
	Authorize(ctx context.Context, accessToken, method string) (*authsvc.UserClaims, error)
}

type service struct {
	tm           authsvc.TokenManager
	protectedRPC map[string][]authsvc.Role
}

func NewService(tm authsvc.TokenManager, protectedRPC map[string][]authsvc.Role) Service {
	return &service{tm, protectedRPC}
}

// Authorize decides whether the user has access to the requested RPC
func (s *service) Authorize(_ context.Context, accessToken, method string) (*authsvc.UserClaims, error) {
	// validate the arguments
	errValid := valid.EmptyError()
	if accessToken == "" {
		errValid.Add("accessToken", "access token is required")
	}
	if method == "" {
		errValid.Add("method", "method is required")
	}
	if !errValid.IsEmpty() {
		return nil, errValid
	}
	// verify the token and extract its claims
	claims, err := s.tm.Verify(accessToken)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", authsvc.ErrInvalidAccessToken, err)
	}
	// convert string role from the claims to authsvc.Role
	role, err := authsvc.StringToRole(claims.Role)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid role: %s: %v", authsvc.ErrInvalidAccessToken, role, err)
	}
	// check whether the requested RPC is protected and if so, whether the user's role has access to it
	roles, ok := s.protectedRPC[method]
	if ok {
		for _, r := range roles {
			if r == role {
				return claims, nil
			}
		}
		return nil, ErrUnauthorized
	}
	return claims, nil
}
