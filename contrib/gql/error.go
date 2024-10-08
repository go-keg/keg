package gql

import (
	"context"
	"errors"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	zlog "github.com/go-keg/keg/contrib/log"
	"github.com/go-keg/keg/contrib/response"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/vektah/gqlparser/v2/gqlerror"
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
	code string
}

type ErrOption func(*ErrOptions)

func WithErrCode(code string) ErrOption {
	return func(opt *ErrOptions) {
		opt.code = code
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
		Message: message,
		Extensions: map[string]any{
			"code": e.code,
		},
		Rule: CustomErrorKey,
	}
}

func ErrorPresenter(logger log.Logger) func(ctx context.Context, err error) *gqlerror.Error {
	return func(ctx context.Context, err error) *gqlerror.Error {
		var gqlErr *gqlerror.Error
		if errors.As(err, &gqlErr) && gqlErr.Err == nil {
			return gqlErr
		}
		code := response.Err2HashCode(err)
		_ = logger.Log(log.LevelError,
			"module", "graphql/errors",
			"traceId", zlog.TraceID(),
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
