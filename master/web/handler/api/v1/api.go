package v1

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	chaosv1 "github.com/slok/ragnarok/api/chaos/v1"
	"github.com/slok/ragnarok/apimachinery/serializer"
	"github.com/slok/ragnarok/log"
	"github.com/slok/ragnarok/master/service/scheduler"
)

// Handler is the handler that has all the required handlers to create the rest api V1
// of the application.
type Handler interface {
	// Debug will create a new debug experiment on the desired node.
	Debug(w http.ResponseWriter, r *http.Request)
	// WriteExperiment will handle the WR operations on an experiment (create & update).
	WriteExperiment(w http.ResponseWriter, r *http.Request)
}

// JSONHandler is the base implementation of Handler using JSON format. Satisfies Handler interface.
type JSONHandler struct {
	scheduler  scheduler.Scheduler
	serializer serializer.Serializer
	logger     log.Logger
}

// NewJSONHandler returns a new api v1 JSON handler.
func NewJSONHandler(scheduler scheduler.Scheduler, logger log.Logger) *JSONHandler {
	return &JSONHandler{
		serializer: serializer.JSONSerializerDefault,
		scheduler:  scheduler,
		logger:     logger,
	}
}

func (j *JSONHandler) setInternalError(w http.ResponseWriter, errorStr string) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "application/json")
	body, _ := json.Marshal(map[string]string{
		"error": errorStr,
	})
	w.Write(body)
}

func (j *JSONHandler) setBadRequest(w http.ResponseWriter, errorStr string) {
	w.WriteHeader(http.StatusBadRequest)
	w.Header().Set("Content-Type", "application/json")
	body, _ := json.Marshal(map[string]string{
		"error": errorStr,
	})
	w.Write(body)
}
func (j *JSONHandler) setOK(w http.ResponseWriter, str string) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	body, _ := json.Marshal(map[string]string{
		"msg": str,
	})
	w.Write(body)
}

func (j *JSONHandler) decodeJSONExperiment(body []byte) (chaosv1.Experiment, error) {
	return chaosv1.NewExperiment(), nil
}

// Debug will create a new experiment on the desired node.
func (j *JSONHandler) Debug(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("Not implemented"))
}

// WriteExperiment will get an experiment in JSON and store an experiment on the repository.
func (j *JSONHandler) WriteExperiment(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		// TODO: Check experiment exists before creating a new one.

		// Deserialize experiment.
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			j.setInternalError(w, err.Error())
			return
		}
		expTmp, err := j.serializer.Decode(b)
		if err != nil {
			j.setBadRequest(w, err.Error())
			return
		}

		exp, ok := expTmp.(*chaosv1.Experiment)
		if !ok {
			j.setBadRequest(w, "decoded object is not an experiment")
			return
		}

		if _, err := j.scheduler.Schedule(exp); err != nil {
			j.setInternalError(w, err.Error())
			return
		}

		j.setOK(w, fmt.Sprintf("experiment %s scheduled", exp.Metadata.ID))
		return
	}
	j.setBadRequest(w, "wrong request")
}
