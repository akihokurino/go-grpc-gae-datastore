package subscriber

import (
	"net/http"

	"gae-go-sample/adapter"
	"gae-go-sample/handler"
)

type Handler func(mux *http.ServeMux)

func NewHandler(
	messageService MessageService,
	checkMaintenance adapter.CheckMaintenance) Handler {

	maintenance := func(server http.Handler) http.Handler {
		return handler.ApplyMiddleware(
			server,
			checkMaintenance)
	}

	return func(mux *http.ServeMux) {
		mux.Handle("/_ah/push-handlers/receive_message",
			maintenance(messageService.SubscribeMessage()))
	}
}
