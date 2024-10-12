package template

import (
	"embed"
	_ "embed"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
)

//go:embed *.tmpl
var fs embed.FS

func Template() entc.Option {
	return templateOption(func(t *gen.Template) (*gen.Template, error) {
		return t.ParseFS(fs, "*.tmpl")
	})
}

// templateOption ensures the template instantiate
// once for config and execute the given Option.
func templateOption(next func(t *gen.Template) (*gen.Template, error)) entc.Option {
	return func(cfg *gen.Config) (err error) {
		tmpl, err := next(gen.NewTemplate("database"))
		if err != nil {
			return err
		}
		cfg.Templates = append(cfg.Templates, tmpl)
		return nil
	}
}
