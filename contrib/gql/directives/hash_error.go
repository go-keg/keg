package directives

import (
	"context"
	"errors"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	errors2 "github.com/go-keg/keg/contrib/errors"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func HashError(ctx context.Context, _ any, next graphql.Resolver) (res any, err error) {
	res, err = next(ctx)
	var gqlErr *gqlerror.Error
	if errors.As(err, &gqlErr) && gqlErr.Err == nil {
		return res, gqlErr
	}
	code := errors2.Err2HashCode(err)
	return res, &gqlerror.Error{
		Err:     err,
		Message: fmt.Sprintf("Unknown error, error code is: %s, if you need assistance, please contact administrator", code),
		Path:    graphql.GetPath(ctx),
	}
}
