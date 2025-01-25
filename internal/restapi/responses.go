package restapi

import (
	"context"
	"encoding/json"
	"net/http"
	"protchain/internal/logging"
	"protchain/internal/tracing"
	"protchain/internal/value"
	"protchain/pkg/function"
)

type ServerResponse struct {
	Err        error           `json:"err"`
	Message    string          `json:"message"`
	Status     string          `json:"status"`
	StatusCode int             `json:"status_code"`
	Context    context.Context `json:"context"`
	Payload    interface{}     `json:"payload"`
}

func respondWithJSONPayload(ctx *tracing.Context, data interface{}, status, message string) *ServerResponse {
	payload, err := json.Marshal(data)
	if err != nil {
		return respondWithError(err, "failed to marshal json payload", value.Error, ctx)
	}

	return &ServerResponse{
		Status:     value.Success,
		StatusCode: function.StatusCode(status),
		Message:    message,
		Payload:    payload,
	}
}

// respondWithError logs an error with zap and parses the error to the ServerResponse
func respondWithError(err error, message, status string, tracingContext *tracing.Context) *ServerResponse {
	logging.Log.Error(err.Error(), logging.GetContext(tracingContext, err)...)

	return &ServerResponse{
		Err:        err,
		Message:    message,
		Status:     status,
		StatusCode: function.StatusCode(status),
	}
}

func writeJSONResponse(w http.ResponseWriter, content []byte, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(statusCode)
	if _, err := w.Write(content); err != nil {
		logging.Log.Error("unable to write json response")
	}
}

// writeErrorResponse writes an error response to the client
func writeErrorResponse(w http.ResponseWriter, err error, status, errMessage string) {
	r := respondWithError(err, errMessage, status, nil)
	response, _ := json.Marshal(r)
	writeJSONResponse(w, response, r.StatusCode)
}
