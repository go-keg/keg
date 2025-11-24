package directives

import (
	"context"
	"fmt"
	"regexp"

	"github.com/99designs/gqlgen/graphql"
	"github.com/go-keg/keg/contrib/gql"
	"github.com/spf13/cast"
)

func Validate(ctx context.Context, obj any, next graphql.Resolver, pattern string, message *string) (res any, err error) {
	path := graphql.GetPathContext(ctx)
	if path != nil && path.Field != nil {
		inputObject, ok := obj.(map[string]any)
		if ok {
			fieldValue := inputObject[*path.Field]
			compile := regexp.MustCompile(pattern)
			if compile.MatchString(cast.ToString(fieldValue)) {
				return next(ctx)
			}
			if message != nil {
				return nil, gql.ValidateError(*message)
			}
			return nil, gql.ValidateError("validate failed")
		}
		return nil, fmt.Errorf("unable to cast obj to map[string]interface{}")
	}
	return nil, fmt.Errorf("path.field is nil")
}
