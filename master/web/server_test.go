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
		{name: "Real routing should call the handlers."},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require := require.New(t)
			assert := assert.New(t)

			// Create simple listener.
			l, err := net.Listen("tcp", "127.0.0.1:0")
			require.NoError(err)

			// Mocks.
			mh := &mhandler.Handler{}
			mh.On("Debug", mock.Anything, mock.Anything).Once().Return()
			mh.On("CreateExperiment", mock.Anything, mock.Anything).Once().Return()

			server := web.NewHTTPServer(web.DefaultHTTPRoutes, mh, l, log.Dummy)
			go func() {
				server.Serve()
			}()

			// API v1
			apiV1DebugURL := fmt.Sprintf("http://%s%s", l.Addr(), web.DefaultHTTPRoutes.APIV1.Debug)
			_, err = http.Get(apiV1DebugURL)
			assert.NoError(err)
			apiV1CreateExperimentURL := fmt.Sprintf("http://%s%s", l.Addr(), web.DefaultHTTPRoutes.APIV1.CreateExperiment)
			_, err = http.Get(apiV1CreateExperimentURL)
			assert.NoError(err)

			// Assert handlders where called.
			mh.AssertExpectations(t)
		})
	}
}
