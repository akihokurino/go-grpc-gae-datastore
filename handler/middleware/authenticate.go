package middleware

import (
	"net/http"

	"github.com/pkg/errors"

	"gae-go-sample/adapter"
	"gae-go-sample/domain"
	"gae-go-sample/handler"
	pb "gae-go-sample/proto/go/pb"
)

const (
	authKey = "Authorization"
)

func NewUserAuthenticate(
	contextProvider handler.ContextProvider,
	fireClient adapter.FirebaseClient,
	logger adapter.CompositeLogger,
	userRepository adapter.UserRepository,
	customerRepository adapter.CustomerRepository,
	clientRepository adapter.ClientRepository,
	companyRepository adapter.CompanyRepository) adapter.UserAuthenticate {
	return func(base http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			meID := domain.UserID("")

			if meID.String() == "" {
				client, err := fireClient.AuthClient(ctx)
				if err != nil {
					serverError(w)
					return
				}

				token := r.Header.Get(authKey)
				if len(token) <= 7 {
					unAuthorizeError(w)
					return
				}

				decoded, err := client.VerifyIDToken(token[7:])
				if err != nil {
					unAuthorizeError(w)
					return
				}

				meID = domain.UserID(decoded.UID)
			}

			logger.Info().With(ctx).Printf("request me id is %s", meID.String())

			me, err := userRepository.Get(ctx, meID)
			if err != nil {
				appErr, ok := errors.Cause(err).(domain.AppError)
				if ok && appErr.Type() == domain.ErrorTypeNotFound {
					newContext, _ := contextProvider.WithAuthUID(ctx, meID)
					base.ServeHTTP(w, r.WithContext(newContext))
					return
				}

				unAuthorizeError(w)
				return
			}

			switch me.Role {
			case pb.User_Role_Customer:
				customer, err := customerRepository.Get(ctx, me.CustomerID())
				if err != nil {
					unAuthorizeError(w)
					return
				}

				if customer.IsDenied() {
					denyError(w)
					return
				}
			case pb.User_Role_Client:
				client, err := clientRepository.Get(ctx, me.ClientID())
				if err != nil {
					unAuthorizeError(w)
					return
				}

				company, err := companyRepository.Get(ctx, client.CompanyID)
				if err != nil {
					unAuthorizeError(w)
					return
				}

				if company.IsBan() {
					banError(w)
					return
				}
			default:
				invalidRoleError(w)
				return
			}

			newContext, _ := contextProvider.WithAuthUID(ctx, meID)
			base.ServeHTTP(w, r.WithContext(newContext))
		})
	}
}

func NewAdminAuthenticate(
	contextProvider handler.ContextProvider,
	fireClient adapter.FirebaseClient,
	fireUserRepository adapter.FireUserRepository,
	logger adapter.CompositeLogger) adapter.AdminAuthenticate {
	return func(base http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			client, err := fireClient.AuthClient(ctx)
			if err != nil {
				serverError(w)
				return
			}

			token := r.Header.Get(authKey)
			if len(token) <= 7 {
				unAuthorizeError(w)
				return
			}

			decoded, err := client.VerifyIDToken(token[7:])
			if err != nil {
				unAuthorizeError(w)
				return
			}

			me, err := fireUserRepository.GetAdmin(ctx, domain.AdminUserID(decoded.UID))
			if err != nil {
				unAuthorizeError(w)
				return
			}

			logger.Info().With(ctx).Printf("request user id is %s", me.UID.String())

			newContext, _ := contextProvider.WithAuthAdminUID(ctx, me.UID)
			base.ServeHTTP(w, r.WithContext(newContext))
			return
		})
	}
}

func NewBatchAuthenticate() adapter.BatchAuthenticate {
	return func(base http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("X-Appengine-Cron") == "true" {
				base.ServeHTTP(w, r)
				return
			}

			if r.Header.Get("X-CaptureHTTP-QueueName") != "" {
				base.ServeHTTP(w, r)
				return
			}

			unAuthorizeError(w)
		})
	}
}
