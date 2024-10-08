package template

import (
	"entgo.io/ent/entc/gen"
)

func Template() *gen.Template {
	return gen.NewTemplate("./gql_custom.tmpl")
}
