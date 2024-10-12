package template

import (
	"embed"
	_ "embed"
	"entgo.io/contrib/entgql"
	"entgo.io/ent/entc/gen"
	"log"
)

//go:embed *.tmpl
var fs embed.FS

func Template() *gen.Template {
	template, err := gen.NewTemplate("gql_custom").Funcs(entgql.TemplateFuncs).ParseFS(fs, "*.tmpl")
	if err != nil {
		log.Fatalf("gql template error: %s\n", err)
	}
	return template
}
