package driver

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"entgo.io/ent/dialect"
	"github.com/google/uuid"
)

type QueryMethod string

const (
	MethodExec           QueryMethod = "Exec"
	MethodExecContext    QueryMethod = "ExecContext"
	MethodQuery          QueryMethod = "Query"
	MethodQueryContext   QueryMethod = "QueryContext"
	MethodTx             QueryMethod = "Tx"
	MethodBeginTx        QueryMethod = "BeginTx"
	MethodTxExec         QueryMethod = "Tx.Exec"
	MethodTxExecContext  QueryMethod = "Tx.ExecContext"
	MethodTxQuery        QueryMethod = "Tx.Query"
	MethodTxQueryContext QueryMethod = "Tx.QueryContext"
	MethodTxCommit       QueryMethod = "Tx.Commit"
	MethodTxRollback     QueryMethod = "Tx.Rollback"
)

type QueryLog struct {
	Method    QueryMethod
	Query     string
	Tx        string
	Args      any
	StartTime time.Time
	Duration  time.Duration
	Err       error
}

type BeforeHookFunc func(ctx context.Context, q QueryLog) context.Context
type AfterHookFunc func(ctx context.Context, q QueryLog)

// DebugDriver is a driver that logs all driver operations.
type DebugDriver struct {
	dialect.Driver // underlying driver.
	before         BeforeHookFunc
	after          AfterHookFunc
}

type DebugOption func(*DebugDriver)

func WithBeforeHook(fn BeforeHookFunc) DebugOption {
	return func(d *DebugDriver) {
		d.before = fn
	}
}

func WithAfterHook(fn AfterHookFunc) DebugOption {
	return func(d *DebugDriver) {
		d.after = fn
	}
}

func Debug(d dialect.Driver, opts ...DebugOption) dialect.Driver {
	drv := &DebugDriver{d, nil, nil}
	for _, opt := range opts {
		opt(drv)
	}
	return drv
}

func (d *DebugDriver) Before(ctx context.Context, log QueryLog) context.Context {
	if d.before != nil {
		ctx = d.before(ctx, log)
	}
	return ctx
}

func (d *DebugDriver) After(ctx context.Context, log QueryLog) {
	if d.after != nil {
		d.after(ctx, log)
	}
}

// Exec logs its params and calls the underlying driver Exec method.
func (d *DebugDriver) Exec(ctx context.Context, query string, args, v any) error {
	q := QueryLog{
		Method:    MethodExec,
		Query:     query,
		Args:      args,
		StartTime: time.Now(),
	}
	ctx = d.Before(ctx, q)
	err := d.Driver.Exec(ctx, query, args, v)
	q.Duration = time.Since(q.StartTime)
	q.Err = err
	d.After(ctx, q)
	return err

}

// ExecContext logs its params and calls the underlying driver ExecContext method if it is supported.
func (d *DebugDriver) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	drv, ok := d.Driver.(interface {
		ExecContext(context.Context, string, ...any) (sql.Result, error)
	})
	if !ok {
		return nil, fmt.Errorf("Driver.ExecContext is not supported")
	}
	q := QueryLog{
		Method:    MethodExecContext,
		Query:     query,
		Args:      args,
		StartTime: time.Now(),
	}
	ctx = d.Before(ctx, q)
	result, err := drv.ExecContext(ctx, query, args...)
	q.Duration = time.Since(q.StartTime)
	q.Err = err
	d.After(ctx, q)
	return result, err
}

// Query logs its params and calls the underlying driver Query method.
func (d *DebugDriver) Query(ctx context.Context, query string, args, v any) error {
	q := QueryLog{
		Method:    MethodQuery,
		Query:     query,
		Args:      args,
		StartTime: time.Now(),
	}
	ctx = d.Before(ctx, q)
	err := d.Driver.Query(ctx, query, args, v)
	q.Duration = time.Since(q.StartTime)
	q.Err = err
	d.After(ctx, q)
	return err
}

// QueryContext logs its params and calls the underlying driver QueryContext method if it is supported.
func (d *DebugDriver) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	drv, ok := d.Driver.(interface {
		QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	})
	if !ok {
		return nil, fmt.Errorf("Driver.QueryContext is not supported")
	}
	q := QueryLog{
		Method:    MethodQueryContext,
		Query:     query,
		Args:      args,
		StartTime: time.Now(),
	}
	ctx = d.Before(ctx, q)
	result, err := drv.QueryContext(ctx, query, args...)
	q.Duration = time.Since(q.StartTime)
	q.Err = err
	d.After(ctx, q)
	return result, err
}

// Tx adds an log-id for the transaction and calls the underlying driver Tx command.
func (d *DebugDriver) Tx(ctx context.Context) (dialect.Tx, error) {
	tx, err := d.Driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	id := uuid.New().String()
	q := QueryLog{
		Method:    MethodTx,
		Tx:        id,
		StartTime: time.Now(),
	}
	ctx = d.Before(ctx, q)
	return &DebugTx{tx, id, d.before, d.after, ctx}, nil
}

// BeginTx adds an log-id for the transaction and calls the underlying driver BeginTx command if it is supported.
func (d *DebugDriver) BeginTx(ctx context.Context, opts *sql.TxOptions) (dialect.Tx, error) {
	drv, ok := d.Driver.(interface {
		BeginTx(context.Context, *sql.TxOptions) (dialect.Tx, error)
	})
	if !ok {
		return nil, fmt.Errorf("Driver.BeginTx is not supported")
	}
	tx, err := drv.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	id := uuid.New().String()
	q := QueryLog{
		Method:    MethodBeginTx,
		Tx:        id,
		StartTime: time.Now(),
	}
	ctx = d.Before(ctx, q)
	return &DebugTx{tx, id, d.before, d.after, ctx}, nil
}

// DebugTx is a transaction implementation that logs all transaction operations.
type DebugTx struct {
	dialect.Tx        // underlying transaction.
	id         string // transaction logging id.
	before     BeforeHookFunc
	after      AfterHookFunc
	ctx        context.Context // underlying transaction context.
}

func (d *DebugTx) Before(ctx context.Context, log QueryLog) context.Context {
	if d.before != nil {
		ctx = d.before(ctx, log)
	}
	return ctx
}

func (d *DebugTx) After(ctx context.Context, log QueryLog) {
	if d.after != nil {
		d.after(ctx, log)
	}
}

// Exec logs its params and calls the underlying transaction Exec method.
func (d *DebugTx) Exec(ctx context.Context, query string, args, v any) error {
	q := QueryLog{
		Method:    MethodTxExec,
		Query:     query,
		Args:      args,
		Tx:        d.id,
		StartTime: time.Now(),
	}
	ctx = d.Before(ctx, q)
	err := d.Tx.Exec(ctx, query, args, v)
	q.Duration = time.Since(q.StartTime)
	q.Err = err
	d.After(ctx, q)
	return err
}

// ExecContext logs its params and calls the underlying transaction ExecContext method if it is supported.
func (d *DebugTx) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	drv, ok := d.Tx.(interface {
		ExecContext(context.Context, string, ...any) (sql.Result, error)
	})
	if !ok {
		return nil, fmt.Errorf("Tx.ExecContext is not supported")
	}
	q := QueryLog{
		Method:    MethodTxExecContext,
		Query:     query,
		Args:      args,
		Tx:        d.id,
		StartTime: time.Now(),
	}
	ctx = d.Before(ctx, q)
	result, err := drv.ExecContext(ctx, query, args...)
	q.Duration = time.Since(q.StartTime)
	q.Err = err
	d.After(ctx, q)
	return result, err
}

// Query logs its params and calls the underlying transaction Query method.
func (d *DebugTx) Query(ctx context.Context, query string, args, v any) error {
	q := QueryLog{
		Method:    MethodTxQuery,
		Query:     query,
		Args:      args,
		Tx:        d.id,
		StartTime: time.Now(),
	}
	ctx = d.Before(ctx, q)
	err := d.Tx.Query(ctx, query, args, v)
	q.Duration = time.Since(q.StartTime)
	q.Err = err
	d.After(ctx, q)
	return err
}

// QueryContext logs its params and calls the underlying transaction QueryContext method if it is supported.
func (d *DebugTx) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	drv, ok := d.Tx.(interface {
		QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	})
	if !ok {
		return nil, fmt.Errorf("Tx.QueryContext is not supported")
	}
	q := QueryLog{
		Method:    MethodTxQueryContext,
		Query:     query,
		Args:      args,
		Tx:        d.id,
		StartTime: time.Now(),
	}
	ctx = d.Before(ctx, q)
	result, err := drv.QueryContext(ctx, query, args...)
	q.Duration = time.Since(q.StartTime)
	q.Err = err
	d.After(ctx, q)
	return result, err
}

// Commit logs this step and calls the underlying transaction Commit method.
func (d *DebugTx) Commit() error {
	d.Before(context.Background(), QueryLog{
		Method:    MethodTxCommit,
		Tx:        d.id,
		StartTime: time.Now(),
	})
	return d.Tx.Commit()
}

// Rollback logs this step and calls the underlying transaction Rollback method.
func (d *DebugTx) Rollback() error {
	d.Before(context.Background(), QueryLog{
		Method:    MethodTxRollback,
		Tx:        d.id,
		StartTime: time.Now(),
	})
	return d.Tx.Rollback()
}
