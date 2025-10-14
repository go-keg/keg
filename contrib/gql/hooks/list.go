package hooks

import (
	"fmt"

	"entgo.io/contrib/entgql"
	"entgo.io/ent/entc/gen"
	"github.com/samber/lo"
	"github.com/vektah/gqlparser/v2/ast"
)

func GenerateList(graph *gen.Graph, schema *ast.Schema) error {
	var queryFields ast.FieldList
	for _, node := range graph.Nodes {
		_, ant, err := gqlTypeFromNode(node)
		if err != nil {
			return err
		}
		if ant.RelayConnection {
			name := lo.CamelCase(node.Name)
			listName := fmt.Sprintf("%sList", node.Name)
			names := &entgql.PaginationNames{
				Order:      fmt.Sprintf("%sOrder", node.Name),
				WhereInput: fmt.Sprintf("%sWhereInput", node.Name),
			}
			schema.AddTypes(&ast.Definition{
				Name:        listName,
				Kind:        ast.Object,
				Description: "A connection to a list of items.",
				Fields: []*ast.FieldDefinition{
					{
						Name:        "nodes",
						Type:        ast.ListType(ast.NonNullNamedType(node.Name, nil), nil),
						Description: "A list of nodes.",
					},
					{
						Name:        "totalCount",
						Type:        ast.NonNullNamedType("Int", nil),
						Description: "Identifies the total count of items in the connection.",
					},
				},
			})
			def := &ast.FieldDefinition{
				Name: name + "List",
				Type: ast.NonNullNamedType(listName, nil),
				Arguments: ast.ArgumentDefinitionList{
					{
						Name:         "offset",
						Type:         ast.NonNullNamedType("Int", nil),
						DefaultValue: &ast.Value{Kind: ast.IntValue, Raw: "0"},
						Description:  "The number of elements to skip from the start of the list.",
					},
					{
						Name:         "limit",
						Type:         ast.NonNullNamedType("Int", nil),
						DefaultValue: &ast.Value{Kind: ast.IntValue, Raw: "10"},
						Description:  "The maximum number of elements to return. This value cannot be negative.",
					},
				},
			}
			if _, ok := schema.Types[node.Name+"Edge"]; ok {
				schema.Types[node.Name+"Edge"] = &ast.Definition{
					Name:        node.Name+"Edge",
					Kind:        ast.Object,
					Description: "An edge in a connection.",
					Fields: []*ast.FieldDefinition{
						{
							Name:        "node",
							Type:        ast.NonNullNamedType(node.Name, nil),
							Description: "The item at the end of the edge.",
						},
						{
							Name:        "cursor",
							Type:        ast.NonNullNamedType("Cursor", nil),
							Description: "A cursor for use in pagination.",
						},
					},
				}
			}
			if _, ok := schema.Types[names.Order]; ok {
				orderT := ast.NamedType(names.Order, nil)
				if ant.MultiOrder {
					orderT = ast.ListType(ast.NonNullNamedType(names.Order, nil), nil)
				}
				def.Arguments = append(def.Arguments, &ast.ArgumentDefinition{
					Name:        "orderBy",
					Type:        orderT,
					Description: fmt.Sprintf("Ordering options for %s returned from the connection.", name),
				})
			}
			if _, ok := schema.Types[names.WhereInput]; ok {
				def.Arguments = append(def.Arguments, &ast.ArgumentDefinition{
					Name:        "where",
					Type:        ast.NamedType(names.WhereInput, nil),
					Description: fmt.Sprintf("Filtering options for %s returned from the connection.", name),
				})
			}
			if ant.QueryField != nil {
				def.Description = ant.QueryField.Description
				def.Directives = buildDirectives(ant.QueryField.Directives)
				queryFields = append(queryFields, def)
			}
		}
	}
	for s, definition := range schema.Types {
		if s == entgql.QueryType {
			definition.Fields = append(definition.Fields, queryFields...)
		}
	}
	return nil
}

func gqlTypeFromNode(t *gen.Type) (gqlType string, ant *entgql.Annotation, err error) {
	if ant, err = annotation(t.Annotations); err != nil {
		return
	}
	gqlType = t.Name
	if ant.Type != "" {
		gqlType = ant.Type
	}
	return
}

// annotation extracts the entgql.Annotation or returns its empty value.
func annotation(ants gen.Annotations) (*entgql.Annotation, error) {
	ant := &entgql.Annotation{}
	if ants != nil && ants[ant.Name()] != nil {
		if err := ant.Decode(ants[ant.Name()]); err != nil {
			return nil, err
		}
	}
	return ant, nil
}

func buildDirectives(directives []entgql.Directive) ast.DirectiveList {
	list := make(ast.DirectiveList, 0, len(directives))
	for _, d := range directives {
		list = append(list, &ast.Directive{
			Name:      d.Name,
			Arguments: d.Arguments,
		})
	}
	return list
}
