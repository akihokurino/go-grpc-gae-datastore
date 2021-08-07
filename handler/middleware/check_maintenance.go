package middleware

import (
	"net/http"

	"gae-go-recruiting-server/adapter"
)

func NewCheckMaintenance(logger adapter.CompositeLogger, switchProvider adapter.SwitchProvider) adapter.CheckMaintenance {
	return func(base http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			if !switchProvider.IsMaintenance {
				base.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			maintenanceError(w)
			return
		})
	}
}
