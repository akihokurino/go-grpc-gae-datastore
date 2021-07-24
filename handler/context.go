package handler

import (
	"context"
	"net/http"

	"gae-go-sample/domain"
)

const requestStoreKey = "__request_store_key__"
const responseStoreKey = "__response_store_key__"
const authUidStoreKey = "__auth_uid_store_key__"
const authAdminUidStoreKey = "__auth_admin_uid_store_key__"

type ContextProvider interface {
	WithAuthUID(ctx context.Context, uid domain.UserID) (context.Context, error)
	WithAuthAdminUID(ctx context.Context, uid domain.AdminUserID) (context.Context, error)
	WithHTTPRequest(ctx context.Context, h http.Request) (context.Context, error)
	WithHTTPResponseWriter(ctx context.Context, w http.ResponseWriter) (context.Context, error)
	AuthUID(ctx context.Context) (domain.UserID, bool)
	MustAuthUID(ctx context.Context) domain.UserID
	AuthAdminUID(ctx context.Context) (domain.AdminUserID, bool)
	MustAuthAdminUID(ctx context.Context) domain.AdminUserID
	HttpRequest(ctx context.Context) (http.Request, bool)
	HttpResponseWriter(ctx context.Context) (http.ResponseWriter, bool)
}

type contextProvider struct {
}

func NewContextProvider() ContextProvider {
	return &contextProvider{}
}

func (u *contextProvider) WithAuthUID(ctx context.Context, uid domain.UserID) (context.Context, error) {
	return context.WithValue(ctx, authUidStoreKey, uid), nil
}

func (u *contextProvider) WithAuthAdminUID(ctx context.Context, uid domain.AdminUserID) (context.Context, error) {
	return context.WithValue(ctx, authAdminUidStoreKey, uid), nil
}

func (u *contextProvider) WithHTTPRequest(ctx context.Context, h http.Request) (context.Context, error) {
	return context.WithValue(ctx, requestStoreKey, h), nil
}

func (u *contextProvider) WithHTTPResponseWriter(ctx context.Context, w http.ResponseWriter) (context.Context, error) {
	return context.WithValue(ctx, responseStoreKey, w), nil
}

func (u *contextProvider) AuthUID(ctx context.Context) (domain.UserID, bool) {
	uid, ok := ctx.Value(authUidStoreKey).(domain.UserID)
	return uid, ok
}

func (u *contextProvider) MustAuthUID(ctx context.Context) domain.UserID {
	uid, ok := ctx.Value(authUidStoreKey).(domain.UserID)
	if !ok {
		panic(domain.NewError(domain.ErrorTypeInternal, "server error"))
	}
	return uid
}

func (u *contextProvider) AuthAdminUID(ctx context.Context) (domain.AdminUserID, bool) {
	uid, ok := ctx.Value(authAdminUidStoreKey).(domain.AdminUserID)
	return uid, ok
}

func (u *contextProvider) MustAuthAdminUID(ctx context.Context) domain.AdminUserID {
	uid, ok := ctx.Value(authAdminUidStoreKey).(domain.AdminUserID)
	if !ok {
		panic(domain.NewError(domain.ErrorTypeInternal, "server error"))
	}
	return uid
}

func (u *contextProvider) HttpRequest(ctx context.Context) (http.Request, bool) {
	h, ok := ctx.Value(requestStoreKey).(http.Request)
	return h, ok
}

func (u *contextProvider) HttpResponseWriter(ctx context.Context) (http.ResponseWriter, bool) {
	w, ok := ctx.Value(responseStoreKey).(http.ResponseWriter)
	return w, ok
}
