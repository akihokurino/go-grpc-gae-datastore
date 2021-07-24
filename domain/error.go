package domain

import (
	"github.com/pkg/errors"
)

type ErrorType int

func (t ErrorType) Code() int {
	return int(t)
}

const (
	// 一般的な内部エラー
	ErrorTypeInternal ErrorType = 500
	// 一般的なクライアントリクエスト起因エラー
	ErrorTypeBadRequest = 400
	// 一般的な認証エラー
	ErrorTypeUnauthorized = 401
	// 一般的な権限エラー
	ErrorTypeForbidden = 403
	// 一般的な存在しないリソース参照エラー
	ErrorTypeNotFound = 404
	// 一般的なコンフリクトエラー
	ErrorTypeAlreadyExists = 409
)

var (
	ErrBadRequest              = NewError(ErrorTypeBadRequest, "invalid params")
	ErrQueryIsTooLong          = NewError(ErrorTypeBadRequest, "query is too long")
	ErrContractFileIsNotExists = NewError(ErrorTypeBadRequest, "contract file is not exists")
	ErrEmailAlreadyExists      = NewError(ErrorTypeAlreadyExists, "email already exists")
	ErrUserAlreadyExists       = NewError(ErrorTypeAlreadyExists, "user already exists")
	ErrEntryAlreadyExists      = NewError(ErrorTypeAlreadyExists, "entry already exists")
	ErrContractAlreadyExists   = NewError(ErrorTypeAlreadyExists, "contract already exists")
	ErrCustomerIsNotActive     = NewError(ErrorTypeForbidden, "customer is not active")
	ErrInvalidClient           = NewError(ErrorTypeForbidden, "invalid client")
	ErrUserIsNotMember         = NewError(ErrorTypeForbidden, "user is not member")
	ErrInvalidUserRole         = NewError(ErrorTypeForbidden, "invalid user role")
	ErrInvalidMessageRoomUser  = NewError(ErrorTypeForbidden, "invalid message room user")
	ErrProjectIsNotOpen        = NewError(ErrorTypeForbidden, "project is not open")
	ErrApplyClientNotAccepted  = NewError(ErrorTypeForbidden, "apply client not accepted")
	ErrCustomerDidNotEntry     = NewError(ErrorTypeForbidden, "customer did not entry")
	ErrContractAlreadyCanceled = NewError(ErrorTypeForbidden, "contract already canceled")
	ErrContractAlreadyAccepted = NewError(ErrorTypeForbidden, "contract already accepted")
	ErrContractInProgress      = NewError(ErrorTypeForbidden, "contract in progress")
	ErrForbiddenClientRole     = NewError(ErrorTypeForbidden, "forbidden client role")
	ErrNoSuchEntity            = NewError(ErrorTypeNotFound, "no such entity")
	ErrServerError             = NewError(ErrorTypeInternal, "internal server expectError")
)

type AppError interface {
	error
	Type() ErrorType
}

type appError struct {
	errType ErrorType
	message string
}

func NewError(errType ErrorType, message string) error {
	return &appError{
		errType: errType,
		message: message,
	}
}

func (err *appError) Error() string {
	return err.message
}

func (err *appError) Type() ErrorType {
	return err.errType
}

func IsNoSuchEntityErr(err error) bool {
	appErr, ok := errors.Cause(err).(AppError)
	if !ok {
		return false
	}

	return appErr.Type() == ErrorTypeNotFound
}

func NewBadRequestError(msg string) error {
	return NewError(ErrorTypeBadRequest, msg)
}