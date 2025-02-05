package gql

import (
	"context"
	"reflect"

	"github.com/99designs/gqlgen/graphql"
	"github.com/samber/lo"
)

const (
	FieldCount = "count"
	FieldItems = "items"
)

type pageOption struct {
	defaultSize int
	maxSize     int
	maxItems    int
}

type PageOption func(*pageOption)

func WithDefaultSize(defaultSize int) PageOption {
	return func(opt *pageOption) {
		opt.defaultSize = defaultSize
	}
}

func WithMaxSize(maxSize int) PageOption {
	return func(opt *pageOption) {
		opt.maxSize = maxSize
	}
}

func WithMaxItems(maxItems int) PageOption {
	return func(opt *pageOption) {
		opt.maxItems = maxItems
	}
}

func OffsetLimit(page *int, size *int, opts ...PageOption) (offset int, limit int) {
	v := pageOption{
		defaultSize: 10,
		maxSize:     1000,
		maxItems:    0,
	}
	for _, opt := range opts {
		opt(&v)
	}
	limit = v.defaultSize
	if size != nil && *size > 0 {
		limit = *size
	}
	if v.maxSize != 0 && limit > v.maxSize {
		limit = v.maxSize
	}
	if page != nil && *page > 1 {
		offset = (*page - 1) * limit
	}
	if v.maxItems != 0 && offset+limit > v.maxItems {
		limit = v.maxItems - offset
		if limit < 0 {
			limit = 0
		}
	}
	return
}

// ContainsField 包含字段
func ContainsField(ctx context.Context, field string) bool {
	return lo.Contains(graphql.CollectAllFields(ctx), field)
}

// ContainsCountField 包含 count 字段
func ContainsCountField(ctx context.Context) bool {
	return ContainsField(ctx, FieldCount)
}

// ContainsItemsField 包含 items 字段
func ContainsItemsField(ctx context.Context) bool {
	return ContainsField(ctx, FieldItems)
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
