package errors

import (
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

// json is the json-iterator instance compatible with the standard library.
var json = jsoniter.ConfigCompatibleWithStandardLibrary

// errorResponse is the JSON body for HTTP error responses.
type errorResponse struct {
	Message string `json:"message"`
}

// WriteError writes an HTTP error response with the given status code and error.
func WriteError(w http.ResponseWriter, statusCode int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	_ = json.NewEncoder(w).Encode(errorResponse{Message: err.Error()})
}
