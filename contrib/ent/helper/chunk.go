package helper

import (
	"context"

	"entgo.io/ent/dialect/sql"
)

type chunkImp[Q, T any] interface {
	Offset(offset int) Q
	Limit(limit int) Q
	All(ctx context.Context) ([]*T, error)
	Count(ctx context.Context) (int, error)
	Where(...func(*sql.Selector))
}

// Chunk 分批处理数据
// example:
//
//	client, err := ent.Open("mysql", os.Getenv("ACCOUNT_DB_DSN"))
//	if err != nil {
//		panic(err)
//	}
//	ctx := context.Background()
//	err = enthelper.Chunk(ctx, client.Account.Query(), 100, func(batchIndex int, items []*ent.Account) error {
//		for i, item := range items {
//			fmt.Println(batch, i, item.Email)
//		}
//		return nil
//	})
func Chunk[Q chunkImp[Q, T], T any](ctx context.Context, query Q, chunk int, fn func(batchIndex int, items []*T) error) error {
	count, err := query.Count(ctx)
	if err != nil {
		return err
	}
	for i := 0; i < count; i += chunk {
		items, err := query.Offset(i).Limit(chunk).All(ctx)
		if err != nil {
			return err
		}
		err = fn(i/chunk, items)
		if err != nil {
			return err
		}
	}
	return nil
}
