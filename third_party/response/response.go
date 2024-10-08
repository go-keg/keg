package response

//go:generate protoc -I. --go_out=paths=source_relative:. response.proto
import (
	"errors"
	"fmt"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

const (
	UnknownHttpCode = 500
	UnknownCode     = 100000
	UnknownReason   = "Unknown"

	SuccessHttpCode = 200
	SuccessCode     = 0
	SuccessReason   = "Success"
	SuccessMessage  = ""

	// SupportPackageIsVersion1 this constant should not be referenced by any other code.
	SupportPackageIsVersion1 = true
)

func (x *Response) Error() string {
	return fmt.Sprintf("error: code = %d reason = %s message = %s", x.GetCode(), x.GetReason(), x.GetMessage())
}

// GRPCStatus returns the Status represented by se.
func (x *Response) GRPCStatus() *status.Status {
	s, _ := status.New(ToGRPCCode(int(x.GetCode())), x.GetMessage()).
		WithDetails(&Response{
			HttpCode: x.GetHttpCode(),
			Code:     x.GetCode(),
			Reason:   x.GetReason(),
			Metadata: x.GetMetadata(),
		})
	return s
}

// Is matches each error in the chain with the target value.
func (x *Response) Is(err error) bool {
	if se := new(Response); errors.As(err, &se) {
		return se.Code == x.Code && se.Reason == x.Reason
	}
	return false
}

// WithMetadata with an MD formed by the mapping of key, value.
func (x *Response) WithMetadata(md map[string]string) *Response {
	err := proto.Clone(x).(*Response)
	err.Metadata = md
	return err
}

// New returns an error object for the code, message.
func New(httpCode, code int, reason, message string) *Response {
	return &Response{
		HttpCode: int32(httpCode),
		Code:     int32(code),
		Message:  message,
		Reason:   reason,
	}
}

// NewErrorf returns an error object for the code, message and error info.
func NewErrorf(httpCode, code int, reason, format string, a ...interface{}) *Response {
	return New(httpCode, code, reason, fmt.Sprintf(format, a...))
}

// Newf New(code fmt.Sprintf(format, a...))
func Newf(httpCode, code int, reason, format string, a ...interface{}) *Response {
	return New(httpCode, code, reason, fmt.Sprintf(format, a...))
}

// HttpCode returns the http code for an error.
// It supports wrapped errors.
func HttpCode(err error) int {
	if err == nil {
		return SuccessHttpCode

	}
	return int(FromError(err).HttpCode)
}

// Code returns the code for an error.
// It supports wrapped errors.
func Code(err error) int {
	if err == nil {
		return SuccessCode

	}
	return int(FromError(err).Code)
}

// Reason returns the reason for a particular error.
// It supports wrapped errors.
func Reason(err error) string {
	if err == nil {
		return UnknownReason
	}
	return FromError(err).Reason
}

// FromError try to convert an error to *Error.
// It supports wrapped errors.
func FromError(err error) *Response {
	if err == nil {
		return nil
	}
	if se := new(Response); errors.As(err, &se) {
		return se
	}
	gs, ok := status.FromError(err)
	if ok {
		ret := New(
			FromGRPCCode(gs.Code()),
			UnknownCode,
			UnknownReason,
			gs.Message(),
		)
		for _, detail := range gs.Details() {
			switch d := detail.(type) {
			case *Response:
				ret.HttpCode = d.HttpCode
				ret.Code = d.Code
				ret.Reason = d.Reason
				return ret.WithMetadata(d.Metadata)
			}
		}
		return ret
	}
	return New(UnknownHttpCode, UnknownCode, UnknownReason, err.Error())
}

func (x *Response) GetResponse(data any) map[string]any {
	resp := map[string]any{
		"code":   x.GetCode(),
		"reason": x.GetReason(),
		"msg":    x.GetMessage(),
	}

	if data != nil {
		resp["data"] = data
	}

	if x.Metadata != nil {
		resp["metadata"] = x.Metadata
	}
	return resp
}
