package gql

import (
	"context"
	"errors"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/go-keg/keg/contrib/errs"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"go.opentelemetry.io/otel/trace"
)

const CustomErrorKey string = "custom"

// ErrDeprecated 方法已弃用
var ErrDeprecated = Error("deprecated", WithErrCode("Deprecated"))

// ErrDisabled 已禁用
var ErrDisabled = Error("disabled", WithErrCode("Disabled"))

// ErrNoPermission 无权限
var ErrNoPermission = Error("no permission", WithErrCode("NoPermission"))

// ErrAccessDenied 拒绝访问，无角色或身份
var ErrAccessDenied = Error("access denied", WithErrCode("AccessDenied"))

// ErrUnauthorized 未授权
var ErrUnauthorized = Error("unauthorized", WithErrCode("Unauthorized"))

// ErrNotFound 数据不存在
var ErrNotFound = Error("data not found", WithErrCode("NotFound"))

type ErrOptions struct {
	code       string
	extensions map[string]any
}

func (r ErrOptions) Extensions() map[string]any {
	if r.extensions == nil {
		return map[string]any{"code": r.code}
	}
	r.extensions["code"] = r.code
	return r.extensions
}

type ErrOption func(*ErrOptions)

func WithErrCode(code string) ErrOption {
	return func(opt *ErrOptions) {
		opt.code = code
	}
}

func WithExtensions(data map[string]any) ErrOption {
	return func(opt *ErrOptions) {
		opt.extensions = data
	}
}

func Error(message string, opts ...ErrOption) *gqlerror.Error {
	e := &ErrOptions{
		code: "ErrUndefined",
	}
	for _, opt := range opts {
		opt(e)
	}
	return &gqlerror.Error{
		Message:    message,
		Extensions: e.Extensions(),
		Rule:       CustomErrorKey,
	}
}

func ValidateError(message string) *gqlerror.Error {
	return Error(message, WithErrCode("ErrValidate"))
}

func ErrorPresenter(logger log.Logger) func(ctx context.Context, err error) *gqlerror.Error {
	return func(ctx context.Context, err error) *gqlerror.Error {
		var gqlErr *gqlerror.Error
		if errors.As(err, &gqlErr) && gqlErr.Err == nil {
			return gqlErr
		}
		code := errs.HashCode(err)
		_ = logger.Log(log.LevelError,
			"module", "graphql/errors",
			"traceId", trace.SpanContextFromContext(ctx).TraceID().String(),
			"errCode", code,
			"err", err,
		)
		return &gqlerror.Error{
			Err:     err,
			Message: fmt.Sprintf("Unknown error, error code is: %s, if you need assistance, please contact administrator", code),
			Path:    graphql.GetPath(ctx),
		}
	}
}
