package rest

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/status"
)

type errBody struct {
	Err string `json:"error,omitempty"`
}

func customHTTPError(ctx context.Context, _ *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter,
	_ *http.Request, err error) {
	const fallback = `{"error": "failed to marshal error message"}`

	errStatus := status.Convert(err)
	w.Header().Set("Content-type", marshaler.ContentType(errStatus.Proto()))
	w.WriteHeader(runtime.HTTPStatusFromCode(status.Code(err)))
	jErr := json.NewEncoder(w).Encode(errBody{
		Err: errStatus.Message(),
	})

	if jErr != nil {
		w.Write([]byte(fallback))
	}
}
