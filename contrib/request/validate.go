package request

import (
	"context"

	"github.com/go-kratos/kratos/v2/middleware"
)

type validator interface {
	Validate() error
}

// Validator is a validator middleware.
func Validator(handler middleware.Handler) middleware.Handler {
	return func(ctx context.Context, req any) (reply any, err error) {
		if v, ok := req.(validator); ok {
			if err = v.Validate(); err != nil {
				return nil, err
			}
		}
		return handler(ctx, req)
	}
}
