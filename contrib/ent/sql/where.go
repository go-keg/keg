package sql

import (
	"entgo.io/ent/dialect/sql"
	"github.com/samber/lo"
)

type Where struct {
	ws []*sql.Predicate
}

func (w *Where) Append(p ...*sql.Predicate) {
	s := lo.Filter(p, func(item *sql.Predicate, index int) bool {
		return item != nil
	})
	if len(s) > 0 {
		w.ws = append(w.ws, s...)
	}
}

func (w *Where) AppendOr(p ...*sql.Predicate) {
	s := lo.Filter(p, func(item *sql.Predicate, index int) bool {
		return item != nil
	})
	if len(s) > 0 {
		w.ws = append(w.ws, sql.Or(s...))
	}
}

func (w *Where) Predicate() *sql.Predicate {
	if len(w.ws) > 0 {
		return sql.And(w.ws...)
	}
	return sql.ExprP("1=1")
}
