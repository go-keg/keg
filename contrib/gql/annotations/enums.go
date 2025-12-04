package annotations

import (
	"entgo.io/ent/entc/gen"
	"github.com/samber/lo"
	"github.com/vektah/gqlparser/v2/ast"
)

const EnumName = "enums"

type Enums map[string]string

func (a Enums) Name() string {
	return EnumName
}

func EnumsGQLSchemaHook(graph *gen.Graph, schema *ast.Schema) error {
	for _, node := range graph.Nodes {
		for _, field := range node.Fields {
			for k, v := range field.Annotations {
				if k == EnumName {
					if enums, ok := v.(map[string]any); ok {
						if enum, ok := schema.Types[node.Name+lo.PascalCase(field.Name)]; ok {
							if enum.Kind == ast.Enum {
								for _, value := range enum.EnumValues {
									if item, ok := enums[value.Name]; ok {
										value.Description = item.(string)
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return nil
}
