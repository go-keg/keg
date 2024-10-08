package mixin

import (
	"time"

	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
)

// -------------------------------------------------
// Mixin definition

// TimeMixin implements the ent.Mixin for sharing
// time fields with package schemas.
type TimeMixin struct {
	// We embed the `mixin.Schema` to avoid
	// implementing the rest of the methods.
	mixin.Schema
}

func (TimeMixin) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("created_at"),
	}
}

func (TimeMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Time("created_at").
			Optional().
			Immutable().
			Default(time.Now).
			Annotations(entgql.Skip(entgql.SkipMutationCreateInput|entgql.SkipMutationUpdateInput), entgql.OrderField("created_at")),
		field.Time("updated_at").
			Optional().
			Default(time.Now).
			UpdateDefault(time.Now).
			Annotations(entgql.Skip(entgql.SkipMutationCreateInput|entgql.SkipMutationUpdateInput), entgql.OrderField("updated_at")),
	}
}
