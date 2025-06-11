package hooks

import (
	atlas "ariga.io/atlas/sql/schema"
	"entgo.io/ent/dialect/sql/schema"
)

func FilterTables(tables ...string) schema.DiffHook {
	return func(next schema.Differ) schema.Differ {
		return schema.DiffFunc(func(current, desired *atlas.Schema) ([]atlas.Change, error) {
			changes, err := next.Diff(current, desired)
			if err != nil {
				return nil, err
			}
			var targetTables = make(map[string]bool)
			for _, table := range tables {
				targetTables[table] = true
			}
			return filterChanges(changes, targetTables), nil
		})
	}
}

func filterChanges(changes []atlas.Change, targetTables map[string]bool) (filtered []atlas.Change) {
	for _, change := range changes {
		switch c := change.(type) {
		case *atlas.AddTable:
			if targetTables[c.T.Name] {
				filtered = append(filtered, c)
			}
		case *atlas.ModifyTable:
			if targetTables[c.T.Name] {
				filtered = append(filtered, c)
			}
		case *atlas.DropTable:
			continue
		}
	}
	return
}
