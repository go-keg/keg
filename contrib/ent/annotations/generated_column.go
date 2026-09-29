package annotations

import "entgo.io/ent/schema"

type GeneratedColumn struct {
	Expression string `json:"expression"`
	Stored     bool   `json:"stored"`
}

func (GeneratedColumn) Name() string {
	return "GeneratedColumn"
}

func (g GeneratedColumn) Map() map[string]any {
	return map[string]any{
		"expression": g.Expression,
		"stored":     g.Stored,
	}
}

var _ schema.Annotation = GeneratedColumn{}
