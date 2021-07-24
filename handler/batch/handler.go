package batch

import (
	"net/http"

	"gae-go-sample/adapter"
	"gae-go-sample/handler"
)

type Handler func(mux *http.ServeMux)

func NewHandler(
	projectService ProjectService,
	checkMaintenance adapter.CheckMaintenance,
	batchAuth adapter.BatchAuthenticate,
	captureHTTP adapter.CaptureHTTP) Handler {

	auth := func(server http.Handler) http.Handler {
		return handler.ApplyMiddleware(
			server,
			checkMaintenance,
			batchAuth,
			captureHTTP)
	}

	return func(mux *http.ServeMux) {
		mux.Handle("/batch/support_project_no_entry", auth(projectService.SupportNoEntry()))
		mux.Handle("/batch/support_project_no_message", auth(projectService.SupportNoMessage()))
	}
}
