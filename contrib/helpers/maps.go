package helpers

import (
	"encoding/json"
)

func MapMerge[K comparable, V any](item map[K]V, maps ...map[K]V) map[K]V {
	for _, m := range maps {
		for k, v := range m {
			if _, ok := item[k]; !ok {
				item[k] = v
			}
		}
	}
	return item
}

func ToStringMap(v any) map[string]any {
	jsonStr, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	var m map[string]any
	err = json.Unmarshal(jsonStr, &m)
	if err != nil {
		return nil
	}
	return m
}

func ToStruct(m map[string]any, v any) error {
	marshal, err := json.Marshal(m)
	if err != nil {
		return err
	}
	err = json.Unmarshal(marshal, v)
	if err != nil {
		return err
	}
	return nil
}

func ToMapE(v any) (map[string]any, error) {
	jsonStr, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var m map[string]any
	err = json.Unmarshal(jsonStr, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func MapGet[K comparable, T any](m map[K]T, key K, defaultValue T) T {
	if v, ok := m[key]; ok {
		return v
	}
	return defaultValue
}

func StructToMap(v any) map[string]any {
	m := make(map[string]any)
	marshal, _ := json.Marshal(v)
	_ = json.Unmarshal(marshal, &m)
	return m
}
