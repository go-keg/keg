package errors

import (
	"crypto/md5" //nolint:gosec
	"fmt"
	"strings"
)

var replacer = strings.NewReplacer(" ", "0", "O", "0", "I", "1")

func Err2HashCode(err error) string {
	msg := err.Error()
	h := md5.Sum([]byte(msg)) //nolint:gosec
	code := strings.ToUpper(fmt.Sprintf("%x", h)[0:4])
	replacer.Replace(code)
	return code
}
