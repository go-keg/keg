package template

import "entgo.io/ent/entc"

func Template() entc.Option {
	return entc.TemplateDir("./database.tmpl")
}
