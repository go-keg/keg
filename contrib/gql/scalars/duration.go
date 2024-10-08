package scalars

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/spf13/cast"
)

func MarshalDuration(t time.Duration) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = io.WriteString(w, t.String())
	})
}

func UnmarshalDuration(v any) (time.Duration, error) {
	switch v := v.(type) {
	case string, int, uint, int64, uint64, json.Number:
		return cast.ToDurationE(v)
	default:
		return 0, fmt.Errorf("%T is not an time.Duration", v)
	}
}
