package rest

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

// logMiddleware logs incoming requests
func (s *Server) logMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func(start time.Time) {
			s.log.Info("REST request",
				zap.String("method", r.Method),
				zap.String("protocol", r.Proto),
				zap.String("path", r.URL.Path),
				zap.String("remote_addr", r.RemoteAddr),
				zap.String("user_agent", r.UserAgent()),
				zap.Float64("took_ms", float64(time.Since(start).Nanoseconds())/1000000),
			)
		}(time.Now())
		next.ServeHTTP(w, r)
	}
}
