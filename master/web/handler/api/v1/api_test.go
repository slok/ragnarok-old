package v1_test

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/slok/ragnarok/log"
	webapiv1 "github.com/slok/ragnarok/master/web/handler/api/v1"
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

			h := webapiv1.NewJSONHandler(log.Dummy)

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
			name:      "GET request should return an ok response",
			reqURL:    "http://valhalla.odin/api/v1/experiment",
			reqBody:   "",
			reqMethod: "GET",
			expCode:   200,
			expBody:   "Not implemented",
		},
		{
			name:      "POST request should return an ok response",
			reqURL:    "http://valhalla.odin/api/v1/experiment",
			reqBody:   "",
			reqMethod: "POST",
			expCode:   200,
			expBody:   "Not implemented",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			h := webapiv1.NewJSONHandler(log.Dummy)

			b := bytes.NewBufferString(test.reqBody)
			req := httptest.NewRequest(test.reqMethod, test.reqURL, b)
			w := httptest.NewRecorder()

			h.CreateExperiment(w, req)
			assert.Equal(test.expCode, w.Code)
			assert.Equal(test.expBody, w.Body.String())
		})
	}
}
