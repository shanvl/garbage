package rest

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
)

// errBody is used in error json encoding: "error": "error message", "fields": {"field1": "err1", "field2": "err2"}
type errBody struct {
	Err    string            `json:"error,omitempty"`
	Fields map[string]string `json:"fields,omitempty"`
}

// customHTTPError is used by REST gateway to transform a gRPC error to the convenient json message which looks like:
// "error": "error message", "fields": {"field1": "err1", "field2": "err2"}
func customHTTPError(_ context.Context, _ *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter,
	_ *http.Request, err error) {
	const fallback = `{"error": "failed to marshal error message"}`

	// convert err to gRPC error
	st := status.Convert(err)
	// populate "fields" with error details and "error" with the error's message
	errBody := errBody{Err: st.Message(), Fields: map[string]string{}}
	for _, detail := range st.Details() {
		if d, ok := detail.(*errdetails.BadRequest); ok {
			for _, violation := range d.GetFieldViolations() {
				errBody.Fields[violation.GetField()] = violation.GetDescription()
			}
		}
	}
	// set content-type header
	w.Header().Set("Content-type", marshaler.ContentType(st.Proto()))
	// set error code
	w.WriteHeader(runtime.HTTPStatusFromCode(status.Code(err)))
	// encode the message to json
	jErr := json.NewEncoder(w).Encode(errBody)
	// use the fallback on error
	if jErr != nil {
		w.Write([]byte(fallback))
	}
}
