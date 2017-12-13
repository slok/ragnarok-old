package web_test

import (
	"fmt"
	"net"
	"net/http"
	"testing"

	"github.com/slok/ragnarok/log"
	"github.com/slok/ragnarok/master/web"
	mhandler "github.com/slok/ragnarok/mocks/master/web/handler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestServerServe(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "Real routing should call the registered handlers."},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resourcePath := "/api/test/v1/serve"

			require := require.New(t)
			assert := assert.New(t)

			// Create simple listener.
			l, err := net.Listen("tcp", "127.0.0.1:0")
			require.NoError(err)

			// Mocks.
			mrh := &mhandler.ResourceHandler{}
			mrh.On("GetRoute").Return(resourcePath)
			mrh.On("List", mock.Anything, mock.Anything, mock.Anything).Return()

			server := web.NewHTTPServer(l, log.Dummy)
			server.HandleResource(mrh)
			go func() {
				server.Serve()
			}()

			// API test
			apiURL := fmt.Sprintf("http://%s%s", l.Addr(), resourcePath)
			_, err = http.Get(apiURL)
			assert.NoError(err)
		})
	}
}
