package grpc

import (
	"errors"
	"fmt"

	"github.com/shanvl/garbage/internal/authsvc"
	"github.com/shanvl/garbage/internal/authsvc/authoriz"
	"github.com/shanvl/garbage/pkg/valid"
	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// grpc specific parsing error which occurs when a provided proto timestamp can't be transformed to an entity of
// time.Time
var ErrInvalidTimestamp = errors.New("invalid timestamp")

// handle error transforms a svc's error into appropriate grpc error. It also logs all unrecognized errors
func (s *Server) handleError(err error) error {
	var validErr *valid.ErrValidation
	switch {
	case errors.As(err, &validErr):
		return errWithDetails(codes.InvalidArgument, validErr.Error(), validErr.Fields())
	case errors.Is(err, authsvc.ErrInactiveUser):
		fallthrough
	case errors.Is(err, authsvc.ErrInvalidActivationToken):
		fallthrough
	case errors.Is(err, authsvc.ErrInvalidRefreshToken):
		fallthrough
	case errors.Is(err, ErrInvalidTimestamp):
		fallthrough
	case errors.Is(err, authsvc.ErrUnknownClient):
		fallthrough
	case errors.Is(err, authsvc.ErrUnknownRole):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, authsvc.ErrUnknownUser):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, authsvc.ErrInvalidPassword):
		fallthrough
	case errors.Is(err, authsvc.ErrInvalidAccessToken):
		return status.Error(codes.Unauthenticated, err.Error())
	case errors.Is(err, authoriz.ErrUnauthorized):
		return status.Error(codes.PermissionDenied, err.Error())
	default:
		s.log.Error("internal error", zap.Error(err))
		return status.Error(codes.Internal, "internal svc error")
	}
}

// errWithDetails takes a map[string]string and appends it as the details to grpc error
func errWithDetails(code codes.Code, message string, details map[string]string) error {
	grpcErr := status.New(code, message)
	br := &errdetails.BadRequest{}
	for field, desc := range details {
		v := &errdetails.BadRequest_FieldViolation{
			Field:       field,
			Description: desc,
		}
		br.FieldViolations = append(br.FieldViolations, v)
	}
	e, err := grpcErr.WithDetails(br)
	if err != nil {
		// there should be no error under normal circumstances, so it's better to
		// panic to figure out what's happened instead of silent error passing
		panic(fmt.Sprintf("unexpected error attaching metadata: %v", err))
	}
	return e.Err()
}
