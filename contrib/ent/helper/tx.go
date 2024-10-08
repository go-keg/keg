package helper

import (
	"context"
	"errors"
)

type clientImp[T txImp] interface {
	Tx(ctx context.Context) (T, error)
}

type txImp interface {
	Rollback() error
	Commit() error
}

func WithTx[Q clientImp[T], T txImp](ctx context.Context, client Q, fn func(ctx context.Context, tx T) error) error {
	tx, err := client.Tx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if v := recover(); v != nil {
			_ = tx.Rollback()
			panic(v)
		}
	}()
	if err = fn(ctx, tx); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			err = errors.Join(err, rollbackErr)
		}
		return err
	}
	return tx.Commit()
}
