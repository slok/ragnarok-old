package webapi_test

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/slok/ragnarok/api"
	"github.com/slok/ragnarok/client/repository/webapi"
	"github.com/slok/ragnarok/log"
	mserializer "github.com/slok/ragnarok/mocks/apimachinery/serializer"
	testapi "github.com/slok/ragnarok/test/api"
)

var (
	testObjType     = api.TypeMeta{Kind: "webapi", Version: "test/v2"}
	testObjTypePath = "/api/test/v2/webapi"
)

// newMockServer is a handy util to create mocked server.
func newMockServer(wantErr bool, expMethod, expBody, expPath string, t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Assertions
		assert.Equal(t, expPath, r.URL.Path)
		b, _ := ioutil.ReadAll(r.Body)
		assert.Equal(t, expBody, string(b))
		assert.Equal(t, expMethod, r.Method)

		// Return
		if wantErr {
			http.Error(w, "wanted error", http.StatusNotFound)
		}
		w.WriteHeader(http.StatusOK)
	}))
}

func TestClientCreate(t *testing.T) {
	tests := []struct {
		name      string
		serverErr bool
		expErr    bool
	}{
		{
			name:      "Creating a resource shouldn't return an error.",
			serverErr: false,
			expErr:    false,
		},
		{
			name:      "If there is an error on the serve it should return an error.",
			serverErr: true,
			expErr:    true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require := require.New(t)
			assert := assert.New(t)

			obj := &testapi.TestObj{Version: "test/v2", Kind: "webapi", ID: "test1"}
			expPath := testObjTypePath
			expObjStr := "mockedStr"
			expMethod := "POST"

			// Mocks
			testServer := newMockServer(test.serverErr, expMethod, expObjStr, expPath, t)
			ms := &mserializer.Serializer{}
			ms.On("Encode", obj, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
				w := args.Get(1).(io.Writer)
				w.Write([]byte(expObjStr))
			})
			ms.On("Decode", mock.Anything).Return(obj, nil)

			// Create the client.
			cli := http.DefaultClient
			c, err := webapi.NewClient(testServer.URL, cli, testObjType, ms, log.Dummy)
			require.NoError(err)

			// Execute & test.
			gotObj, err := c.Create(obj)
			if test.expErr {
				assert.Error(err)
			} else if assert.NoError(err) {
				assert.Equal(obj, gotObj)
			}
		})
	}
}
