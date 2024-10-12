package template

import (
	_ "embed"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
)

//go:embed database.tmpl
var temp []byte

func Template() entc.Option {
	return templateOption(func(t *gen.Template) (*gen.Template, error) {
		return t.Parse(string(temp))
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
