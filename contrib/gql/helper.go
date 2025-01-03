package gql

import (
	"context"
	"reflect"

	"github.com/99designs/gqlgen/graphql"
	"github.com/samber/lo"
)

const (
	fieldCount = "count"
	fieldItems = "items"
)

func OffsetLimit(page *int, size *int) (offset int, limit int) {
	if size != nil && *size > 0 {
		limit = *size
	} else {
		limit = 10
	}
	if page != nil && *page > 1 {
		offset = (*page - 1) * limit
	}
	return
}

// ContainsField 包含字段
func ContainsField(ctx context.Context, field string) bool {
	return lo.Contains(graphql.CollectAllFields(ctx), field)
}

// ContainsCountField 包含 count 字段
func ContainsCountField(ctx context.Context) bool {
	return ContainsField(ctx, fieldCount)
}

// ContainsItemsField 包含 items 字段
func ContainsItemsField(ctx context.Context) bool {
	return ContainsField(ctx, fieldItems)
}

// SomeFields 包含任意字段
func SomeFields(ctx context.Context, fields ...string) bool {
	return lo.Some(graphql.CollectAllFields(ctx), fields)
}

// EveryFields 每个字段都包含
func EveryFields(ctx context.Context, fields ...string) bool {
	return lo.Every(graphql.CollectAllFields(ctx), fields)
}

func SomeWhere(items ...any) bool {
	for _, item := range items {
		if !reflect.ValueOf(item).IsNil() {
			return true
		}
	}
	return false
}
