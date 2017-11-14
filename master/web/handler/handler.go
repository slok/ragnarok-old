package handler

import (
	apiv1 "github.com/slok/ragnarok/master/web/handler/api/v1"
)

// Handler is the handler that will serve the server.
type Handler interface {
	apiv1.Handler
}
