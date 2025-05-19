package filter

import (
	"entgo.io/ent/dialect/sql"
	"github.com/samber/lo"
)

type Filter struct {
	ws []*sql.Predicate
}

func (w *Filter) Append(p ...*sql.Predicate) {
	s := lo.Filter(p, func(item *sql.Predicate, _ int) bool {
		return item != nil
	})
	if len(s) > 0 {
		w.ws = append(w.ws, s...)
	}
}

func (w *Filter) AppendOr(p ...*sql.Predicate) {
	s := lo.Filter(p, func(item *sql.Predicate, _ int) bool {
		return item != nil
	})
	if len(s) > 0 {
		w.ws = append(w.ws, sql.Or(s...))
	}
}

func (w *Filter) Predicate() *sql.Predicate {
	if len(w.ws) > 0 {
		return sql.And(w.ws...)
	}
	return sql.ExprP("1=1")
}
