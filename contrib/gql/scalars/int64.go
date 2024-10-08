package scalars

import (
	"encoding/json"
	"fmt"
	"io"

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
