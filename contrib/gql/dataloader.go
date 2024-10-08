package gql

import (
	"context"
	"github.com/graph-gophers/dataloader"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

type LoaderFunc func(ctx context.Context, keys dataloader.Keys) (map[dataloader.Key]any, error)

func BatchFunc(fn LoaderFunc) dataloader.BatchFunc {
	return func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		records, err := fn(ctx, keys)
		if err != nil {
			return []*dataloader.Result{{Data: nil, Error: err}}
		}
		results := make([]*dataloader.Result, len(keys))
		for i, key := range keys {
			if v, ok := records[key]; ok {
				results[i] = &dataloader.Result{Data: v, Error: nil}
			} else {
				results[i] = &dataloader.Result{Data: nil, Error: nil}
			}
		}
		return results
	}
}

func ToInts(keys dataloader.Keys) []int {
	return lo.Map(keys, func(item dataloader.Key, index int) int {
		return cast.ToInt(item.String())
	})
}

func ToAnySlice(keys dataloader.Keys) []any {
	return lo.Map(keys, func(item dataloader.Key, index int) any {
		return cast.ToInt(item.String())
	})
}

func ToStringKey(id any) dataloader.Key {
	return dataloader.StringKey(cast.ToString(id))
}
