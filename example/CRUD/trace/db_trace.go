package trace

import (
	"context"
	"database/sql"

	"go.opentelemetry.io/otel/attribute"
	db "go.opentelemetry.io/otel/example/CRUD/db/sqlc"
)

var (
	queryKey = attribute.Key("ex.com/db/query")
)

type DBTXTrace struct {
	db db.DBTX
}

func NewDBTXTrace(db db.DBTX) db.DBTX {
	return &DBTXTrace{db: db}
}

func (tracer *DBTXTrace) ExecContext(ctx context.Context, str string, opts ...interface{}) (result sql.Result, err error) {
	tr := traceProvider.Tracer("example/curd/database-transaction")
	ctx, span := tr.Start(ctx, "ExecContext")
	defer span.End()
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetAttributes(queryKey.String(str))
		}
	}()
	result, err = tracer.db.ExecContext(ctx, str, opts...)
	return
}

func (tracer *DBTXTrace) PrepareContext(ctx context.Context, str string) (stmt *sql.Stmt, err error) {
	tr := traceProvider.Tracer("example/curd/database-transaction")
	ctx, span := tr.Start(ctx, "PrepareContext")
	defer span.End()
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetAttributes(queryKey.String(str))
		}
	}()

	stmt, err = tracer.db.PrepareContext(ctx, str)
	return
}

func (tracer *DBTXTrace) QueryContext(ctx context.Context, str string, opts ...interface{}) (rows *sql.Rows, err error) {
	tr := traceProvider.Tracer("example/curd/database-transaction")
	ctx, span := tr.Start(ctx, "QueryContext")
	defer span.End()
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetAttributes(queryKey.String(str))
		}
	}()

	rows, err = tracer.db.QueryContext(ctx, str, opts...)
	return
}

func (tracer *DBTXTrace) QueryRowContext(ctx context.Context, str string, opts ...interface{}) *sql.Row {
	tr := traceProvider.Tracer("example/curd/database-transaction")
	ctx, span := tr.Start(ctx, "QueryRowContext")
	defer span.End()
	return tracer.db.QueryRowContext(ctx, str, opts...)
}
