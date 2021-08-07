package middleware

import (
	"net/http"

	"gae-go-recruiting-server/adapter"
	"gae-go-recruiting-server/handler"
)

func NewCaptureHTTP(contextProvider handler.ContextProvider) adapter.CaptureHTTP {
	return func(base http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			ctx, err := contextProvider.WithHTTPRequest(ctx, *r)
			if err != nil {
				serverError(w)
				return
			}

			ctx, err = contextProvider.WithHTTPResponseWriter(ctx, w)
			if err != nil {
				serverError(w)
				return
			}

			base.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
