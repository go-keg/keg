package annotations

import (
	"entgo.io/ent/entc/gen"
	"github.com/vektah/gqlparser/v2/ast"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
	"unicode"
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
						if enum, ok := schema.Types[node.Name+Studly(field.Name)]; ok {
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

// Studly returns the String instance in studly case.
func Studly(s string) string {
	words := fieldsFunc(s, func(r rune) bool {
		return r == '_' || r == ' ' || r == '-' || r == ',' || r == '.'
	}, func(r rune) bool {
		return unicode.IsUpper(r)
	})

	casesTitle := cases.Title(language.Und)
	var studlyWords []string
	for _, word := range words {
		studlyWords = append(studlyWords, casesTitle.String(word))
	}

	return strings.Join(studlyWords, "")
}

// fieldsFunc splits the input string into words with preservation, following the rules defined by
// the provided functions f and preserveFunc.
func fieldsFunc(s string, f func(rune) bool, preserveFunc ...func(rune) bool) []string {
	var fields []string
	var currentField strings.Builder

	shouldPreserve := func(r rune) bool {
		for _, preserveFn := range preserveFunc {
			if preserveFn(r) {
				return true
			}
		}
		return false
	}

	for _, r := range s {
		if f(r) {
			if currentField.Len() > 0 {
				fields = append(fields, currentField.String())
				currentField.Reset()
			}
		} else if shouldPreserve(r) {
			if currentField.Len() > 0 {
				fields = append(fields, currentField.String())
				currentField.Reset()
			}
			currentField.WriteRune(r)
		} else {
			currentField.WriteRune(r)
		}
	}

	if currentField.Len() > 0 {
		fields = append(fields, currentField.String())
	}

	return fields
}
