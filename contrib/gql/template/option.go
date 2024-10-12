package template

import (
	_ "embed"
	"entgo.io/contrib/entgql"
	"entgo.io/ent/entc/gen"
)

//go:embed gql_custom.tmpl
var temp []byte

func Template() *gen.Template {
	template, err := gen.NewTemplate("gql_custom").Parse(string(temp))
	if err != nil {
		return nil
	}
	return template.Funcs(entgql.TemplateFuncs)
}
