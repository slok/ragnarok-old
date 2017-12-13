package v1_test

import (
	"bytes"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/slok/ragnarok/api"
	clusterv1 "github.com/slok/ragnarok/api/cluster/v1"
	"github.com/slok/ragnarok/apimachinery/serializer"
	webapiclusterv1 "github.com/slok/ragnarok/master/web/handler/api/cluster/v1"
	mcliclusterv1 "github.com/slok/ragnarok/mocks/client/api/cluster/v1"
)

func TestNodeHandlerCreate(t *testing.T) {
	tests := []struct {
		name      string
		reqBody   string
		createErr bool
		expNode   *clusterv1.Node
		expCode   int
		expBody   string
	}{
		{
			name:    "Request to create a new node with an empty json should return an error.",
			reqBody: "",
			expCode: 400,
			expBody: `{"error":"unknown type of object: unexpected end of JSON input"}`,
		},
		{
			name:    "Request to create a new node with an wrong type should return an error.",
			reqBody: `{"kind":"experiment","version":"chaos/v1","metadata":{"id":"exp-001","name":"first experiment","description":" first experiment is the first experiment :|"},"spec":{"selector":{"az":"eu-west-1a","kind":"master"},"template":{"spec":{"timeout":300000000000,"attacks":[{"attack1":{"size":524288000}},{"attack2":null},{"attack3":{"pace":"10m","quantity":10,"rest":"30s","target":"myTarget"}}]}}},"status":{"failureIDs":["node1","node3","node4"],"creation":"2012-11-01T22:08:41Z"}}`,
			expCode: 400,
			expBody: `{"error":"decoded object is not a node"}`,
		},
		{
			name:      "Request to create a new node should return an error if creating the node, this should returns an error.",
			reqBody:   `{"kind":"node","version":"cluster/v1","metadata":{"id":"testNode1","labels":{"id":"testNode1","kind":"node"},"annotations":{"name":"my node"}}}`,
			createErr: true,
			expNode:   nil,
			expCode:   500,
			expBody:   `{"error":"error"}`,
		},
		{
			name:    "Request to create a new node should return the created node.",
			reqBody: `{"kind":"node","version":"cluster/v1","metadata":{"id":"testNode1","labels":{"id":"testNode1","kind":"node"},"annotations":{"name":"my node"}}}`,
			expNode: &clusterv1.Node{
				TypeMeta: clusterv1.NodeTypeMeta,
				Metadata: api.ObjectMeta{
					ID: "testNode1",
					Labels: map[string]string{
						"kind": "node",
						"id":   "testNode1",
					},
					Annotations: map[string]string{
						"name": "my node",
					},
				},
			},
			expCode: 200,
			expBody: `{"kind":"node","version":"cluster/v1","metadata":{"id":"testNode1","labels":{"id":"testNode1","kind":"node"},"annotations":{"name":"my node"}},"spec":{},"status":{"creation":"0001-01-01T00:00:00Z"}}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			var createErr error
			if test.createErr {
				createErr = errors.New("error")
			}

			// Mocks.
			mcv1 := &mcliclusterv1.NodeClientInterface{}
			mcv1.On("Create", mock.Anything).Return(test.expNode, createErr)

			nh := webapiclusterv1.NewNodeHandler(serializer.DefaultSerializer, mcv1)
			b := bytes.NewBufferString(test.reqBody)
			r := httptest.NewRequest("POST", "http://test", b)
			w := httptest.NewRecorder()

			nh.Create(w, r)
			assert.Equal(test.expCode, w.Code)
			assert.Equal(test.expBody, strings.TrimSuffix(w.Body.String(), "\n"))
		})
	}
}

func TestNodeHandlerUpdate(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		reqBody string
		expCode int
		expBody string
	}{
		{
			name:    "Request to update a node should return not implemented.",
			reqBody: "",
			id:      "",
			expCode: 500,
			expBody: `{"error":"not implemented"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			// Mocks.
			mcv1 := &mcliclusterv1.NodeClientInterface{}

			nh := webapiclusterv1.NewNodeHandler(serializer.DefaultSerializer, mcv1)
			b := bytes.NewBufferString(test.reqBody)
			r := httptest.NewRequest("POST", "http://test", b)
			w := httptest.NewRecorder()

			nh.Update(w, r, test.id)
			assert.Equal(test.expCode, w.Code)
			assert.Equal(test.expBody, strings.TrimSuffix(w.Body.String(), "\n"))
		})
	}
}

func TestNodeHandlerDelete(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		reqBody string
		expCode int
		expBody string
	}{
		{
			name:    "Request to delete a node should return not implemented.",
			reqBody: "",
			id:      "",
			expCode: 500,
			expBody: `{"error":"not implemented"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			// Mocks.
			mcv1 := &mcliclusterv1.NodeClientInterface{}

			nh := webapiclusterv1.NewNodeHandler(serializer.DefaultSerializer, mcv1)
			b := bytes.NewBufferString(test.reqBody)
			r := httptest.NewRequest("POST", "http://test", b)
			w := httptest.NewRecorder()

			nh.Delete(w, r, test.id)
			assert.Equal(test.expCode, w.Code)
			assert.Equal(test.expBody, strings.TrimSuffix(w.Body.String(), "\n"))
		})
	}
}

func TestNodeHandlerGet(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		reqBody string
		expCode int
		expBody string
	}{
		{
			name:    "Request to get a node should return not implemented.",
			reqBody: "",
			id:      "",
			expCode: 500,
			expBody: `{"error":"not implemented"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			// Mocks.
			mcv1 := &mcliclusterv1.NodeClientInterface{}

			nh := webapiclusterv1.NewNodeHandler(serializer.DefaultSerializer, mcv1)
			b := bytes.NewBufferString(test.reqBody)
			r := httptest.NewRequest("POST", "http://test", b)
			w := httptest.NewRecorder()

			nh.Get(w, r, test.id)
			assert.Equal(test.expCode, w.Code)
			assert.Equal(test.expBody, strings.TrimSuffix(w.Body.String(), "\n"))
		})
	}
}

func TestNodeHandlerList(t *testing.T) {
	tests := []struct {
		name    string
		opts    map[string]string
		reqBody string
		expCode int
		expBody string
	}{
		{
			name:    "Request to list nodes should return not implemented.",
			reqBody: "",
			opts:    map[string]string{},
			expCode: 500,
			expBody: `{"error":"not implemented"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			// Mocks.
			mcv1 := &mcliclusterv1.NodeClientInterface{}

			nh := webapiclusterv1.NewNodeHandler(serializer.DefaultSerializer, mcv1)
			b := bytes.NewBufferString(test.reqBody)
			r := httptest.NewRequest("POST", "http://test", b)
			w := httptest.NewRecorder()

			nh.List(w, r, test.opts)
			assert.Equal(test.expCode, w.Code)
			assert.Equal(test.expBody, strings.TrimSuffix(w.Body.String(), "\n"))
		})
	}
}

func TestNodeHandlerWatch(t *testing.T) {
	tests := []struct {
		name    string
		opts    map[string]string
		reqBody string
		expCode int
		expBody string
	}{
		{
			name:    "Request to watch nodes should return not implemented.",
			reqBody: "",
			opts:    map[string]string{},
			expCode: 500,
			expBody: `{"error":"not implemented"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			// Mocks.
			mcv1 := &mcliclusterv1.NodeClientInterface{}

			nh := webapiclusterv1.NewNodeHandler(serializer.DefaultSerializer, mcv1)
			b := bytes.NewBufferString(test.reqBody)
			r := httptest.NewRequest("POST", "http://test", b)
			w := httptest.NewRecorder()

			nh.Watch(w, r, test.opts)
			assert.Equal(test.expCode, w.Code)
			assert.Equal(test.expBody, strings.TrimSuffix(w.Body.String(), "\n"))
		})
	}
}

func TestNodeHandlerRoute(t *testing.T) {
	assert := assert.New(t)
	expRoute := "/api/cluster/v1/node"

	// Mocks.
	mcv1 := &mcliclusterv1.NodeClientInterface{}

	nh := webapiclusterv1.NewNodeHandler(serializer.DefaultSerializer, mcv1)
	assert.Equal(expRoute, nh.GetRoute())
}
