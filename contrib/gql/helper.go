package gql

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/go-keg/keg/contrib/helpers"
	"github.com/samber/lo"
	"github.com/vektah/gqlparser/v2/ast"
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

// QueryFieldsHash generates a unique MD5 hash based on the selected GraphQL fields from the given query context.
func QueryFieldsHash(ctx context.Context) string {
	fields := graphql.CollectFieldsCtx(ctx, nil)
	var content []string
	for _, field := range fields {
		content = append(content, fieldNames(field.Field)...)
	}
	sort.Slice(content, func(i, j int) bool {
		return content[i] < content[j]
	})
	return helpers.MD5(strings.Join(content, ","))
}

func fieldNames(field *ast.Field) []string {
	var names []string
	if len(field.SelectionSet) > 0 {
		var fs []string
		for _, selection := range field.SelectionSet {
			v := fieldNames(selection.(*ast.Field))
			fs = append(fs, v...)
		}
		names = append(names, fmt.Sprintf("%s{%s}", field.Name, strings.Join(fs, ",")))
	} else {
		names = append(names, field.Name)
	}
	return names
}
