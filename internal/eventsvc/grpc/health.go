package grpc

import (
	"context"

	healthv1pb "github.com/shanvl/garbage/api/health/v1/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Check is used for health checks
func (s *Server) Check(_ context.Context, _ *healthv1pb.HealthCheckRequest) (*healthv1pb.HealthCheckResponse, error) {
	return &healthv1pb.HealthCheckResponse{Status: healthv1pb.HealthCheckResponse_SERVING}, nil
}

// Watch is used for stream health checks. Not implemented but required by gRPC Health Checking Protocol
func (s *Server) Watch(_ *healthv1pb.HealthCheckRequest, _ healthv1pb.Health_WatchServer) error {
	return status.Errorf(codes.Unimplemented, "no stream health check")
}
