package gql

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"strconv"

	"github.com/graph-gophers/dataloader"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

type IntKey int

func (k IntKey) String() string {
	return strconv.Itoa(int(k))
}

func (k IntKey) Raw() any {
	return int(k)
}

type StringsKey []byte

func NewStringsKey(values ...string) dataloader.Key {
	var buf bytes.Buffer
	for _, v := range values {
		err := binary.Write(&buf, binary.BigEndian, int32(len(v)))
		if err != nil {
			return nil
		}
		buf.WriteString(v)
	}
	return StringsKey(buf.Bytes())
}

func (r StringsKey) String() string {
	return string(r)
}

func (r StringsKey) Raw() any {
	return []byte(r)
}

func (r StringsKey) Split() ([]string, error) {
	var result []string
	buf := bytes.NewReader(r)
	for buf.Len() > 0 {
		var length int32
		if err := binary.Read(buf, binary.BigEndian, &length); err != nil {
			return nil, err
		}
		content := make([]byte, length)
		if _, err := buf.Read(content); err != nil {
			return nil, err
		}
		result = append(result, string(content))
	}
	return result, nil
}

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
	return lo.Map(keys, func(item dataloader.Key, _ int) int {
		return cast.ToInt(item.String())
	})
}

func ToInt64s(keys dataloader.Keys) []int64 {
	return lo.Map(keys, func(item dataloader.Key, _ int) int64 {
		return cast.ToInt64(item.String())
	})
}

func ToStrings(keys dataloader.Keys) []string {
	return lo.Map(keys, func(item dataloader.Key, _ int) string {
		return item.String()
	})
}

func ToAnySlice(keys dataloader.Keys) []any {
	return lo.Map(keys, func(item dataloader.Key, _ int) any {
		return item.Raw()
	})
}

func ToStringKey(id any) dataloader.Key {
	return dataloader.StringKey(cast.ToString(id))
}

func FillDefault(keys dataloader.Keys, result map[dataloader.Key]any, value any) map[dataloader.Key]any {
	for _, key := range keys {
		if _, ok := result[key]; !ok {
			result[key] = value
		}
	}
	return result
}

func FillDefaultByKey(keys dataloader.Keys, result map[dataloader.Key]any, fn func(dataloader.Key) any) map[dataloader.Key]any {
	for _, key := range keys {
		if _, ok := result[key]; !ok {
			result[key] = fn(key)
		}
	}
	return result
}

func LoadManyResult[T any](items []any, err []error) ([]T, error) {
	if err != nil || items == nil {
		return nil, errors.Join(err...)
	}
	return lo.FilterMap(items, func(item any, _ int) (T, bool) {
		var t T
		if item == nil {
			return t, false
		}
		return item.(T), true
	}), nil
}
