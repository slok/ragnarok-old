package v1

import (
	"net/http"

	"github.com/slok/ragnarok/log"
)

// Handler is the handler that has all the required handlers to create the rest api V1
// of the application.
type Handler interface {
	// Debug will create a new debug experiment on the desired node.
	Debug(w http.ResponseWriter, r *http.Request)
	// CreateExperiment will create a new experiment on the master node.
	CreateExperiment(w http.ResponseWriter, r *http.Request)
}

// JSONHandler is the base implementation of Handler using JSON format. Satisfies Handler interface.
type JSONHandler struct {
	logger log.Logger
}

// NewJSONHandler returns a new api v1 JSON handler.
func NewJSONHandler(logger log.Logger) *JSONHandler {
	return &JSONHandler{
		logger: logger,
	}
}

// Debug will create a new experiment on the desired node.
func (j *JSONHandler) Debug(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Not implemented"))
}

// CreateExperiment will get an experiment in JSON and store an experiment on the repository.
func (j *JSONHandler) CreateExperiment(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Not implemented"))
}
