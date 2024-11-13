package directives

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/go-keg/keg/contrib/gql"
)

func Disabled(ctx context.Context, obj any, next graphql.Resolver) (res any, err error) {
	return nil, gql.ErrDeprecated
}
