package sql

import (
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqljson"
)

func In[T any](col string, items []T) *sql.Predicate {
	if len(items) == 1 {
		return sql.EQ(col, items[0])
	} else if len(items) > 1 {
		var result = make([]any, len(items))
		for i, item := range items {
			result[i] = item
		}
		return sql.In(col, result...)
	}
	return nil
}

func NotIn[T any](col string, items []T) *sql.Predicate {
	if len(items) == 1 {
		return sql.NEQ(col, items[0])
	} else if len(items) > 1 {
		var result = make([]any, len(items))
		for i, item := range items {
			result[i] = item
		}
		return sql.NotIn(col, result...)
	}
	return nil
}

func JsonValueContains[T any](col string, items []T) *sql.Predicate {
	if len(items) > 0 {
		var ws []*sql.Predicate
		for _, item := range items {
			ws = append(ws, sqljson.ValueContains(col, item))
		}
		return sql.Or(ws...)
	}
	return nil
}

func EQ[T any](col string, v *T) *sql.Predicate {
	var t *T
	if v != t {
		return sql.EQ(col, v)
	}
	return nil
}

func Contains(col string, v *string) *sql.Predicate {
	var t *string
	if v != t {
		return sql.Contains(col, *v)
	}
	return nil
}

func FulltextMatch(col string, v *string) *sql.Predicate {
	if v != nil {
		return sql.ExprP(fmt.Sprintf("MATCH (%s) AGAINST ('%s')", col, *v))
	}
	return nil
}

func NotNull(col string, b *bool) *sql.Predicate {
	if b == nil {
		return nil
	}
	if *b {
		return sql.NotNull(col)
	}
	return sql.IsNull(col)
}
