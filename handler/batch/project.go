package batch

import (
	"net/http"
	"time"

	"gae-go-recruiting-server/adapter"
)

type ProjectService interface {
	SupportNoEntry() http.Handler
	SupportNoMessage() http.Handler
}

type projectHandler struct {
	projectApplication adapter.ProjectApplication
}

func NewProjectHandler(projectApplication adapter.ProjectApplication) ProjectService {
	return &projectHandler{
		projectApplication: projectApplication,
	}
}

func (h *projectHandler) SupportNoEntry() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		now := time.Now()

		if err := h.projectApplication.SupportNoEntry(ctx, now); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	})
}

func (h *projectHandler) SupportNoMessage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		now := time.Now()

		if err := h.projectApplication.SupportNoMessage(ctx, now); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	})
}
