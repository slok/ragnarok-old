package web_test

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/slok/ragnarok/master/web"
	mhandler "github.com/slok/ragnarok/mocks/master/web/handler"
)

func TestResourceDispatcherDispatchCreate(t *testing.T) {
	tests := []struct {
		name      string
		reqURL    string
		reqBody   string
		reqMethod string
	}{
		{
			name:      "POST request should create a new resource.",
			reqURL:    "http://valhalla.odin/api/test/v1/test",
			reqBody:   "",
			reqMethod: "POST",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Mocks.
			mrh := &mhandler.ResourceHandler{}
			mrh.On("Create", mock.Anything, mock.Anything)

			b := bytes.NewBufferString(test.reqBody)
			req := httptest.NewRequest(test.reqMethod, test.reqURL, b)
			w := httptest.NewRecorder()

			dispatcher := web.NewResourceHandlerDispatcher(mrh)
			dispatcher.Dispatch(w, req)

			mrh.AssertExpectations(t)
		})
	}
}

func TestResourceDispatcherDispatchGet(t *testing.T) {
	tests := []struct {
		name         string
		handlerRoute string
		reqURL       string
		reqBody      string
		reqMethod    string
		expID        string
	}{
		{
			name:         "GET request should retrieve a resource.",
			handlerRoute: "api/test/v1/test",
			reqURL:       "http://valhalla.odin/api/test/v1/test/myid001",
			reqBody:      "",
			reqMethod:    "GET",
			expID:        "myid001",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Mocks.
			mrh := &mhandler.ResourceHandler{}
			mrh.On("GetRoute").Return(test.handlerRoute)
			mrh.On("Get", mock.Anything, mock.Anything, test.expID)

			b := bytes.NewBufferString(test.reqBody)
			req := httptest.NewRequest(test.reqMethod, test.reqURL, b)
			w := httptest.NewRecorder()

			dispatcher := web.NewResourceHandlerDispatcher(mrh)
			dispatcher.Dispatch(w, req)

			mrh.AssertExpectations(t)
		})
	}
}

func TestResourceDispatcherDispatchList(t *testing.T) {
	tests := []struct {
		name         string
		handlerRoute string
		reqURL       string
		reqBody      string
		reqMethod    string
		expOpts      map[string]string
	}{
		{
			name:         "GET request should retrieve a resource.",
			handlerRoute: "api/test/v1/test",
			reqURL:       "http://valhalla.odin/api/test/v1/test/?nolabelSelector=name=test,version=v1",
			reqBody:      "",
			reqMethod:    "GET",
			expOpts:      map[string]string{},
		},
		{
			name:         "GET request should retrieve a resource with options.",
			handlerRoute: "api/test/v1/test",
			reqURL:       "http://valhalla.odin/api/test/v1/test/?labelSelector=name=test,version=v1",
			reqBody:      "",
			reqMethod:    "GET",
			expOpts:      map[string]string{"version": "v1", "name": "test"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Mocks.
			mrh := &mhandler.ResourceHandler{}
			mrh.On("GetRoute").Return(test.handlerRoute)
			mrh.On("List", mock.Anything, mock.Anything, test.expOpts)

			b := bytes.NewBufferString(test.reqBody)
			req := httptest.NewRequest(test.reqMethod, test.reqURL, b)
			w := httptest.NewRecorder()

			dispatcher := web.NewResourceHandlerDispatcher(mrh)
			dispatcher.Dispatch(w, req)

			mrh.AssertExpectations(t)
		})
	}
}

func TestResourceDispatcherDispatchWatch(t *testing.T) {
	tests := []struct {
		name         string
		handlerRoute string
		reqURL       string
		reqBody      string
		reqMethod    string
		expOpts      map[string]string
	}{
		{
			name:         "GET request should retrieve a watcher to list resources.",
			handlerRoute: "api/test/v1/test",
			reqURL:       "http://valhalla.odin/api/test/v1/test/?nolabelSelector=name=test,version=v1&watch=TRUE",
			reqBody:      "",
			reqMethod:    "GET",
			expOpts:      map[string]string{},
		},
		{
			name:         "GET request should retrieve a watcher to list resources with options.",
			handlerRoute: "api/test/v1/test",
			reqURL:       "http://valhalla.odin/api/test/v1/test/?watch=true&labelSelector=name=test,version=v1",
			reqBody:      "",
			reqMethod:    "GET",
			expOpts:      map[string]string{"version": "v1", "name": "test"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Mocks.
			mrh := &mhandler.ResourceHandler{}
			mrh.On("GetRoute").Return(test.handlerRoute)
			mrh.On("Watch", mock.Anything, mock.Anything, test.expOpts)

			b := bytes.NewBufferString(test.reqBody)
			req := httptest.NewRequest(test.reqMethod, test.reqURL, b)
			w := httptest.NewRecorder()

			dispatcher := web.NewResourceHandlerDispatcher(mrh)
			dispatcher.Dispatch(w, req)

			mrh.AssertExpectations(t)
		})
	}
}

func TestResourceDispatcherDispatchDelete(t *testing.T) {
	tests := []struct {
		name         string
		handlerRoute string
		reqURL       string
		reqBody      string
		reqMethod    string
		expID        string
	}{
		{
			name:         "DELETE request should delete a resource.",
			handlerRoute: "api/test/v1/test",
			reqURL:       "http://valhalla.odin/api/test/v1/test/myid001",
			reqBody:      "",
			reqMethod:    "DELETE",
			expID:        "myid001",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Mocks.
			mrh := &mhandler.ResourceHandler{}
			mrh.On("GetRoute").Return(test.handlerRoute)
			mrh.On("Delete", mock.Anything, mock.Anything, test.expID)

			b := bytes.NewBufferString(test.reqBody)
			req := httptest.NewRequest(test.reqMethod, test.reqURL, b)
			w := httptest.NewRecorder()

			dispatcher := web.NewResourceHandlerDispatcher(mrh)
			dispatcher.Dispatch(w, req)

			mrh.AssertExpectations(t)
		})
	}
}

func TestResourceDispatcherDispatchUpdate(t *testing.T) {
	tests := []struct {
		name         string
		handlerRoute string
		reqURL       string
		reqBody      string
		reqMethod    string
		expID        string
	}{
		{
			name:         "DELETE request should delete a resource.",
			handlerRoute: "api/test/v1/test",
			reqURL:       "http://valhalla.odin/api/test/v1/test/myid001",
			reqBody:      "",
			reqMethod:    "PUT",
			expID:        "myid001",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Mocks.
			mrh := &mhandler.ResourceHandler{}
			mrh.On("GetRoute").Return(test.handlerRoute)
			mrh.On("Update", mock.Anything, mock.Anything, test.expID)

			b := bytes.NewBufferString(test.reqBody)
			req := httptest.NewRequest(test.reqMethod, test.reqURL, b)
			w := httptest.NewRecorder()

			dispatcher := web.NewResourceHandlerDispatcher(mrh)
			dispatcher.Dispatch(w, req)

			mrh.AssertExpectations(t)
		})
	}
}
