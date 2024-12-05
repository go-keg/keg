package directives

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/go-keg/keg/contrib/cache"
)

func Cache(ctx context.Context, _ any, next graphql.Resolver, duration string) (res any, err error) {
	field := graphql.GetFieldContext(ctx)
	marshal, err := json.Marshal(field.Args)
	if err != nil {
		return nil, err
	}
	encoded := base64.StdEncoding.EncodeToString(marshal)
	cacheKey := fmt.Sprintf("gqlCache:%s:%s", field.Field.Name, encoded)
	d, err := time.ParseDuration(duration)
	if err != nil {
		return nil, err
	}
	return cache.LocalRemember(cacheKey, d, func() (any, error) {
		return next(ctx)
	})
}
