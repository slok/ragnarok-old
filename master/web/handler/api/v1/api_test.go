package v1_test

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/slok/ragnarok/log"
	webapiv1 "github.com/slok/ragnarok/master/web/handler/api/v1"
	mclichaosv1 "github.com/slok/ragnarok/mocks/client/api/chaos/v1"
)

func TestJSONHandlerDebug(t *testing.T) {
	tests := []struct {
		name      string
		reqURL    string
		reqBody   string
		reqMethod string
		expCode   int
		expBody   string
	}{
		{
			name:      "GET request should return an ok response",
			reqURL:    "http://valhalla.odin/api/v1/debug",
			reqBody:   "",
			reqMethod: "GET",
			expCode:   200,
			expBody:   "Not implemented",
		},
		{
			name:      "POST request should return an ok response",
			reqURL:    "http://valhalla.odin/api/v1/debug",
			reqBody:   "",
			reqMethod: "POST",
			expCode:   200,
			expBody:   "Not implemented",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			// Mocks.
			mce := &mclichaosv1.ExperimentClientInterface{}
			mce.On("Create", mock.Anything).Return(nil, nil)

			h := webapiv1.NewJSONHandler(mce, log.Dummy)

			b := bytes.NewBufferString(test.reqBody)
			req := httptest.NewRequest(test.reqMethod, test.reqURL, b)
			w := httptest.NewRecorder()

			h.Debug(w, req)
			assert.Equal(test.expCode, w.Code)
			assert.Equal(test.expBody, w.Body.String())
		})
	}
}

func TestJSONHandlerCreateExperiment(t *testing.T) {
	tests := []struct {
		name      string
		reqURL    string
		reqBody   string
		reqMethod string
		expCode   int
		expBody   string
	}{
		{
			name:      "GET request should return an error.",
			reqURL:    "http://valhalla.odin/api/v1/experiment",
			reqBody:   "",
			reqMethod: "GET",
			expCode:   400,
			expBody:   `{"error":"wrong request"}`,
		},
		{
			name:      "POST request with wrong body should return an error.",
			reqURL:    "http://valhalla.odin/api/v1/experiment",
			reqBody:   "",
			reqMethod: "POST",
			expCode:   400,
			expBody:   `{"error":"unknown type of object: unexpected end of JSON input"}`,
		},
		{
			name:      "POST request with wrong object should return an error.",
			reqURL:    "http://valhalla.odin/api/v1/experiment",
			reqBody:   `{"kind":"node","version":"cluster/v1","metadata":{"id":"testNode1","master":true},"spec":{"labels":{"id":"testNode1","kind":"node"}},"status":{"state":1,"creation":"2012-11-01T22:08:41Z"}}`,
			reqMethod: "POST",
			expCode:   400,
			expBody:   `{"error":"decoded object is not an experiment"}`,
		},
		{
			name:      "POST request with correct experiment object shouldn't return an error.",
			reqURL:    "http://valhalla.odin/api/v1/experiment",
			reqBody:   `{"kind":"experiment","version":"chaos/v1","metadata":{"id":"exp-001","name":"first experiment","description":" first experiment is the first experiment :|"},"spec":{"selector":{"az":"eu-west-1a","kind":"master"},"template":{"spec":{"timeout":300000000000,"attacks":[{"attack1":{"size":524288000}},{"attack2":null},{"attack3":{"pace":"10m","quantity":10,"rest":"30s","target":"myTarget"}}]}}},"status":{"failureIDs":["node1","node3","node4"],"creation":"2012-11-01T22:08:41Z"}}`,
			reqMethod: "POST",
			expCode:   200,
			expBody:   `{"msg":"experiment exp-001 scheduled"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			// Mocks.
			mce := &mclichaosv1.ExperimentClientInterface{}
			mce.On("Create", mock.Anything).Return(nil, nil)

			h := webapiv1.NewJSONHandler(mce, log.Dummy)

			b := bytes.NewBufferString(test.reqBody)
			req := httptest.NewRequest(test.reqMethod, test.reqURL, b)
			w := httptest.NewRecorder()

			h.WriteExperiment(w, req)
			assert.Equal(test.expCode, w.Code)
			assert.Equal(test.expBody, w.Body.String())
		})
	}
}
