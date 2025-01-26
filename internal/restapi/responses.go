package restapi

import (
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
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
	File       io.Reader       `json:"-"`
	FileName   string          `json:"-"`
	FileType   string          `json:"-"`
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

// writeMultipartResponse writes a multipart response to the client
func writeMultipartResponse(w http.ResponseWriter, jsonData []byte, file io.Reader, fileName string, fileType string) error {
	mw := multipart.NewWriter(w)
	defer mw.Close()

	// Set the Content-Type header
	w.Header().Set("Content-Type", mw.FormDataContentType())

	// Write JSON part
	jsonPart, err := mw.CreatePart(map[string][]string{"Content-Type": {"application/json"}})
	if err != nil {
		return err
	}
	if _, err := jsonPart.Write(jsonData); err != nil {
		return err
	}

	// Write file part
	filePart, err := mw.CreatePart(map[string][]string{
		"Content-Disposition": {`attachment; filename="` + fileName + `"`},
		"Content-Type":        {fileType},
	})
	if err != nil {
		return err
	}
	if _, err := io.Copy(filePart, file); err != nil {
		return err
	}

	return nil
}

// writeErrorResponse writes an error response to the client
func writeErrorResponse(w http.ResponseWriter, err error, status, errMessage string) {
	r := respondWithError(err, errMessage, status, nil)
	response, _ := json.Marshal(r)
	writeJSONResponse(w, response, r.StatusCode)
}
