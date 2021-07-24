package handler

import (
	"context"
	"strconv"

	"gae-go-sample/adapter"

	"gae-go-sample/domain"

	"github.com/pkg/errors"
	"github.com/twitchtv/twirp"
)

func NewErrorConverter(logger adapter.CompositeLogger) adapter.ErrorConverter {
	return func(ctx context.Context, err error) error {
		appErr, ok := errors.Cause(err).(domain.AppError)

		var (
			twirpErrCode  twirp.ErrorCode
			domainErrCode int
			message       string
		)

		if !ok {
			// ドメインエラーでラップされていない未識別のエラーの場合、エラーメッセージが露出するとセキュリティリスクになり得るので隠蔽する
			twirpErrCode = twirp.Internal
			message = "non specified internal error"
			domainErrCode = domain.ErrorTypeInternal.Code()
		} else {
			twirpErrCode, ok = domainErrTypeToTwirpErrCodeMap[appErr.Type()]
			if !ok {
				twirpErrCode = twirp.Internal
			}
			// ドメインエラーのメッセージはレスポンスに露出させる
			message = appErr.Error()
			domainErrCode = appErr.Type().Code()
		}

		if twirpErrCode == twirp.Internal {
			logger.Error().With(ctx).Printf("%+v", err)
		}

		return twirp.NewError(twirpErrCode, message).
			WithMeta("Code", strconv.Itoa(domainErrCode))
	}
}

var domainErrTypeToTwirpErrCodeMap = map[domain.ErrorType]twirp.ErrorCode{
	domain.ErrorTypeInternal:      twirp.Internal,
	domain.ErrorTypeBadRequest:    twirp.InvalidArgument,
	domain.ErrorTypeUnauthorized:  twirp.Unauthenticated,
	domain.ErrorTypeForbidden:     twirp.PermissionDenied,
	domain.ErrorTypeNotFound:      twirp.NotFound,
	domain.ErrorTypeAlreadyExists: twirp.AlreadyExists,
}
