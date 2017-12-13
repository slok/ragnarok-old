package util

import (
	"encoding/json"
	"net/http"
)

// SetJSONInternalError sets an internal HTTP error.
func SetJSONInternalError(w http.ResponseWriter, errorStr string) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "application/json")
	body, _ := json.Marshal(map[string]string{
		"error": errorStr,
	})
	w.Write(body)
}

// SetJSONBadRequest sets a bad request HTTP error.
func SetJSONBadRequest(w http.ResponseWriter, errorStr string) {
	w.WriteHeader(http.StatusBadRequest)
	w.Header().Set("Content-Type", "application/json")
	body, _ := json.Marshal(map[string]string{
		"error": errorStr,
	})
	w.Write(body)
}

// SetJSONNotImplementedError sets a not implemented HTTP response.
func SetJSONNotImplementedError(w http.ResponseWriter) {
	SetJSONInternalError(w, "not implemented")
}

func SetJSONOK(w http.ResponseWriter, b []byte) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}
