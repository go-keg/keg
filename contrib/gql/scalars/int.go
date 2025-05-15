package scalars

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
	"github.com/spf13/cast"
)

func MarshalInt64(t int64) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = io.WriteString(w, fmt.Sprintf(`"%d"`, t))
	})
}

func UnmarshalInt64(v any) (int64, error) {
	switch v := v.(type) {
	case string, int, uint, int64, uint64:
		return cast.ToInt64E(v)
	case json.Number:
		return v.Int64()
	default:
		return 0, fmt.Errorf("%T is not an int64", v)
	}
}

func MarshalInt8(t int8) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = io.WriteString(w, fmt.Sprintf(`"%d"`, t))
	})
}

func UnmarshalInt8(v any) (int8, error) {
	switch v := v.(type) {
	case string, int, uint, int8, uint8:
		return cast.ToInt8E(v)
	default:
		return 0, fmt.Errorf("%T is not an int8", v)
	}
}

func MarshalUint8(t uint8) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = io.WriteString(w, fmt.Sprintf(`"%d"`, t))
	})
}

func UnmarshalUint8(v any) (uint8, error) {
	switch v := v.(type) {
	case string, int, uint, int8, uint8:
		return cast.ToUint8E(v)
	default:
		return 0, fmt.Errorf("%T is not an uint8", v)
	}
}

func MarshalUint(t uint) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = io.WriteString(w, strconv.FormatUint(uint64(t), 10))
	})
}

func UnmarshalUint(v any) (uint, error) {
	switch v := v.(type) {
	case string, int, int64, uint, uint64:
		return cast.ToUintE(v)
	default:
		return 0, fmt.Errorf("%T is not an uint", v)
	}
}
