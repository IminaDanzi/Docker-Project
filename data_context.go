package dat

import (
	"context"

	"github.com/jackc/pgx/v4"

	"github.com/pkg/errors"
)

// Context provides accessors to all database entities
type Context struct {
	ctx         context.Context
	db          pgx.Tx
	RequestBody *RequestBodyMapper
}

//NewDataContext creates a transaction and uses it to provide accessors to all DB entities
func NewDataContext(mainCtx context.Context) (*Context, error) {
	if dbMap == nil {
		return nil, errors.New("database connection is nil")
	}

	txn, err := dbMap.Begin(mainCtx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to start transaction")
	}

	return &Context{
		ctx:         mainCtx,
		db:          txn,
		RequestBody: NewRequestBodyMapper(mainCtx, txn),
	}, nil
}

// Rollback rolls back the transaction in this context
func (c *Context) Rollback() error {
	return c.db.Rollback(c.ctx)
}

// Commit commits the transaction in this context
func (c *Context) Commit() error {
	return c.db.Commit(c.ctx)
}
