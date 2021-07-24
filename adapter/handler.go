package adapter

import (
	"context"
	"net/http"
)

type ErrorConverter func(ctx context.Context, err error) error

type CaptureHTTP func(base http.Handler) http.Handler

type UserAuthenticate func(base http.Handler) http.Handler

type AdminAuthenticate func(base http.Handler) http.Handler

type BatchAuthenticate func(base http.Handler) http.Handler

type CROS func(base http.Handler) http.Handler

type CheckMaintenance func(base http.Handler) http.Handler
