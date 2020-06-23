package grpc

import (
	"errors"

	"github.com/shanvl/garbage/internal/eventssvc"
	"github.com/shanvl/garbage/internal/eventssvc/eventing"
	"github.com/shanvl/garbage/pkg/valid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// handle error transforms a service's error into appropriate grpc error. It also logs all unrecognized errors
func (s *Server) handleError(err error) error {
	var validErr *valid.ErrValidation
	switch {
	case errors.As(err, &validErr):
		return status.Error(codes.InvalidArgument, validErr.Error())
	case errors.Is(err, eventssvc.ErrInvalidClassName):
		fallthrough
	case errors.Is(err, eventssvc.ErrNoClassOnDate):
		fallthrough
	case errors.Is(err, eventssvc.ErrUnknownResource):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, eventssvc.ErrUnknownPupil):
		fallthrough
	case errors.Is(err, eventing.ErrNoEventPupil):
		fallthrough
	case errors.Is(err, eventssvc.ErrUnknownEvent):
		return status.Error(codes.NotFound, err.Error())
	default:
		s.log.Error("internal error", zap.Error(err))
		return status.Error(codes.Internal, "internal service error")
	}
}
