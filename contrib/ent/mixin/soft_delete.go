package mixin

import (
	"context"
	"fmt"
	"time"

	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
)

// SoftDeleteMixin implements the soft delete pattern for schemas.
type SoftDeleteMixin struct {
	mixin.Schema
}

func (SoftDeleteMixin) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("deleted_at"),
	}
}

// Fields of the SoftDeleteMixin.
func (SoftDeleteMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Time("deleted_at").Optional().Annotations(entgql.Skip()),
	}
}

type softDeleteKey struct{}

// SkipSoftDelete returns a new context that skips the soft-delete interceptor/mutators.
func SkipSoftDelete(parent context.Context) context.Context {
	return context.WithValue(parent, softDeleteKey{}, true)
}

// Interceptors of the SoftDeleteMixin.
func (d SoftDeleteMixin) Interceptors() []ent.Interceptor {
	return []ent.Interceptor{
		ent.TraverseFunc(func(ctx context.Context, query ent.Query) error {
			// Skip soft-delete, means include soft-deleted entities.
			if skip, _ := ctx.Value(softDeleteKey{}).(bool); skip {
				return nil
			}
			if q, ok := query.(interface {
				WhereP(...func(*sql.Selector))
			}); ok {
				d.P(q)
			}
			return nil
		}),
	}
}

// Hooks of the SoftDeleteMixin.
func (d SoftDeleteMixin) Hooks() []ent.Hook {
	return []ent.Hook{
		func(next ent.Mutator) ent.Mutator {
			return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
				if !m.Op().Is(ent.OpDeleteOne | ent.OpDelete) {
					return next.Mutate(ctx, m)
				}
				// Skip soft-delete, means delete the entity permanently.
				if skip, _ := ctx.Value(softDeleteKey{}).(bool); skip {
					return next.Mutate(ctx, m)
				}
				mx, ok := m.(interface {
					SetOp(ent.Op)
					Mutate(context.Context, ent.Mutation) (ent.Value, error)
					SetDeletedAt(time.Time)
					WhereP(...func(*sql.Selector))
				})
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				d.P(mx)
				mx.SetOp(ent.OpUpdate)
				mx.SetDeletedAt(time.Now())
				return mx.Mutate(ctx, m)
			})
		},
	}
}

// P adds a storage-level predicate to the queries and mutations.
func (d SoftDeleteMixin) P(w interface{ WhereP(...func(*sql.Selector)) }) {
	w.WhereP(sql.FieldIsNull(d.Fields()[0].Descriptor().Name))
}
