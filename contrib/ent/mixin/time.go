package mixin

import (
	"strings"
	"time"

	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
	"github.com/samber/lo"
)

// -------------------------------------------------
// Mixin definition

// TimeMixin implements the ent.Mixin for sharing
// time fields with package schemas.
type TimeMixin struct {
	// We embed the `mixin.Schema` to avoid
	// implementing the rest of the methods.
	mixin.Schema
	// SortFieldCaseStyle default value is NamingStyleSnakeCase
	SortFieldCaseStyle SortFieldCaseStyle
}

type SortFieldCaseStyle string

const (
	NamingStylePascalCase SortFieldCaseStyle = "PascalCase"
	NamingStyleCamelCase  SortFieldCaseStyle = "camelCase"
	NamingStyleSnakeCase  SortFieldCaseStyle = "snake_case"
	NamingStyleUpperCase  SortFieldCaseStyle = "UPPER_CASE"
)

func (TimeMixin) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("created_at"),
	}
}

func (r TimeMixin) Fields() []ent.Field {
	createdAtOrder := "created_at"
	updatedAtOrder := "updated_at"
	switch r.SortFieldCaseStyle {
	case NamingStylePascalCase:
		createdAtOrder = lo.PascalCase(createdAtOrder)
		updatedAtOrder = lo.PascalCase(updatedAtOrder)
	case NamingStyleCamelCase:
		createdAtOrder = lo.CamelCase(createdAtOrder)
		updatedAtOrder = lo.CamelCase(updatedAtOrder)
	case NamingStyleUpperCase:
		createdAtOrder = strings.ToUpper(createdAtOrder)
		updatedAtOrder = strings.ToUpper(updatedAtOrder)
	}
	return []ent.Field{
		field.Time("created_at").
			Optional().
			Immutable().
			Default(time.Now).
			Annotations(
				entgql.Skip(entgql.SkipMutationCreateInput|entgql.SkipMutationUpdateInput),
				entgql.OrderField(createdAtOrder),
			),
		field.Time("updated_at").
			Optional().
			Default(time.Now).
			UpdateDefault(time.Now).
			Annotations(
				entgql.Skip(entgql.SkipMutationCreateInput|entgql.SkipMutationUpdateInput),
				entgql.OrderField(updatedAtOrder),
			),
	}
}
