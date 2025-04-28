package gql

import (
	"fmt"
	"io"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
)

func MarshalerString(s string) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = w.Write([]byte(strconv.Quote(s)))
	})
}

func UnmarshalerString[T ~string](v any) (T, error) {
	switch v := v.(type) {
	case string:
		return T(v), nil
	default:
		return "", fmt.Errorf("%T is not a string", v)
	}
}

func MarshalerUint8(s uint8) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = w.Write([]byte(strconv.Itoa(int(s))))
	})
}

func UnmarshalerUint8[T ~uint8](v any) (T, error) {
	switch v := v.(type) {
	case uint8:
		return T(v), nil
	default:
		return 0, fmt.Errorf("%T is not a uint8", v)
	}
}
