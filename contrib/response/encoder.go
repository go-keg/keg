package response

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/go-keg/keg/third_party/response"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"
	nethttp "net/http"
	"strings"
)

type ErrorEncoderOption func(resp *response.Response, err error) *response.Response

var replacer = strings.NewReplacer(" ", "0", "O", "0", "I", "1")

func Err2HashCode(err error) string {
	msg := err.Error()
	h := md5.Sum([]byte(msg))
	code := strings.ToUpper(fmt.Sprintf("%x", h)[0:4])
	replacer.Replace(code)
	return code
}

func HashUnknownError(logger log.Logger) ErrorEncoderOption {
	return func(resp *response.Response, err error) *response.Response {
		if resp.GetCode() == response.UnknownCode {
			code := Err2HashCode(err)
			_ = logger.Log(log.LevelError, "code", code, "msg", err)
			resp.Message = fmt.Sprintf("Unknown error, error code is: %s, if you need assistance, please contact the administrator", code)
			resp = resp.WithMetadata(map[string]string{
				"code": code,
			})
		}
		return resp
	}
}

func ErrorEncoder(opts ...ErrorEncoderOption) http.EncodeErrorFunc {
	return func(w nethttp.ResponseWriter, r *nethttp.Request, err error) {
		e := response.FromError(err)
		for _, opt := range opts {
			e = opt(e, err)
		}
		resp := e.GetResponse(nil)
		body, _ := json.Marshal(resp)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(int(e.GetHttpCode()))
		_, _ = w.Write(body)
	}
}

func Encoder() http.EncodeResponseFunc {
	return func(w nethttp.ResponseWriter, r *nethttp.Request, v interface{}) error {
		codec, _ := http.CodecForRequest(r, "Accept")
		data, err := codec.Marshal(v)
		if err != nil {
			return err
		}
		data = []byte(fmt.Sprintf(`{"code":%d,"data":%s,"msg":"%s","reason":"%s"}`,
			response.SuccessCode,
			data,
			response.SuccessMessage,
			response.SuccessReason,
		))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(nethttp.StatusOK)
		_, err = w.Write(data)
		return err
	}
}
