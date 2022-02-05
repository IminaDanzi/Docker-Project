package dat

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"github.com/jackc/pgx/v4"
)

const (
	requestBodyTable           = "requestbody"
	requestBodyTableBodyColumn = "body"
)

// RequestBodyMapper provides CRUD methods for the table requestbody
type RequestBodyMapper struct {
	ctx  context.Context
	txn  pgx.Tx
	psql sq.StatementBuilderType
}

// NewRequestBodyMapper creates a new RequestBodyMapper for transaction
func NewRequestBodyMapper(ctx context.Context, txn pgx.Tx) *RequestBodyMapper {
	return &RequestBodyMapper{
		ctx:  ctx,
		txn:  txn,
		psql: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

// Insert will try to create a row in the table 'requestbody' with the provided input
func (rb *RequestBodyMapper) Insert(input []byte) error {
	sql, params, err := rb.psql.
		Insert(requestBodyTable).
		Columns(requestBodyTableBodyColumn).
		Values(input).
		ToSql()
	if err != nil {
		return err
	}

	rows, err := rb.txn.Query(rb.ctx, sql, params...)
	rows.Close()
	return err
}
