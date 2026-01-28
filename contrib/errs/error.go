package errs

import (
	"errors"
)

type Error struct {
	Code    int            // 业务错误码
	Message string         // 对外错误信息
	Cause   error          // 原始错误（仅日志）
	Meta    map[string]any // 扩展信息（调试/定位）
}

func (e *Error) Error() string {
	return e.Message
}

func New(code int, msg string) *Error {
	return &Error{
		Code:    code,
		Message: msg,
	}
}

func Wrap(err error, code int, msg string) *Error {
	return &Error{
		Code:    code,
		Message: msg,
		Cause:   err,
	}
}

func Is(err error, code int) bool {
	var e *Error
	ok := errors.As(err, &e)
	return ok && e.Code == code
}

func Cause(err error) error {
	var e *Error
	if errors.As(err, &e) {
		return e.Cause
	}
	return err
}

func WithMeta(err error, meta map[string]any) *Error {
	var e *Error
	ok := errors.As(err, &e)
	if !ok {
		return &Error{
			Code:    ErrInternal,
			Message: "系统异常",
			Cause:   err,
			Meta:    meta,
		}
	}
	e.Meta = meta
	return e
}
