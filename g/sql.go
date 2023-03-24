package g

import (
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

type QueryerContext interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

type QueryRowerContext interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type ExecerContext interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type Sqlizer interface {
	ToSql() (query string, args []any, err error)
}

type SelectDoc[T any] interface {
	*T
	ScanPointers() []any
}

func FindManyDocs[T any, S SelectDoc[T]](ctx context.Context, queryer QueryerContext, sqlizer Sqlizer) ([]T, error) {
	query, args, err := sqlizer.ToSql()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	rows, err := queryer.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer rows.Close()
	var result []T
	for rows.Next() {
		doc := new(T)
		err = rows.Scan(S(doc).ScanPointers()...)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		result = append(result, *doc)
	}
	return result, nil
}

func FindManyValues[T any](ctx context.Context, queryer QueryerContext, sqlizer Sqlizer) ([]T, error) {
	query, args, err := sqlizer.ToSql()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	rows, err := queryer.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer rows.Close()
	var result []T
	for rows.Next() {
		var val T
		err = rows.Scan(&val)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		result = append(result, val)
	}
	return result, nil
}

func FindOneDoc[T any, S SelectDoc[T]](ctx context.Context, queryer QueryRowerContext, sqlizer Sqlizer) (T, bool, error) {
	var empty T
	query, args, err := sqlizer.ToSql()
	if err != nil {
		return empty, false, errors.WithStack(err)
	}
	doc := new(T)
	err = queryer.QueryRowContext(ctx, query, args...).Scan(S(doc).ScanPointers()...)
	if err != nil {
		if err == sql.ErrNoRows {
			return empty, false, nil
		}
		return empty, false, errors.WithStack(err)
	}
	return *doc, true, nil
}

func FindOneValue[T any](ctx context.Context, queryer QueryRowerContext, sqlizer Sqlizer) (T, bool, error) {
	var empty T
	query, args, err := sqlizer.ToSql()
	if err != nil {
		return empty, false, errors.WithStack(err)
	}
	var val T
	err = queryer.QueryRowContext(ctx, query, args...).Scan(&val)
	if err != nil {
		if err == sql.ErrNoRows {
			return empty, false, nil
		}
		return empty, false, errors.WithStack(err)
	}
	return val, true, nil
}

type UUID string

func InsertValuesReturningId(ctx context.Context, queryer QueryRowerContext, sqlizer Sqlizer) (UUID, error) {
	query, args, err := sqlizer.ToSql()
	if err != nil {
		return Empty[UUID](), errors.WithStack(err)
	}
	query = fmt.Sprintf("%s RETURNING `id`", query)
	var id UUID
	err = queryer.QueryRowContext(ctx, query, args...).Scan(&id)
	if err != nil {
		return Empty[UUID](), errors.WithStack(err)
	}
	return id, nil
}

func InsertValues(ctx context.Context, execer ExecerContext, sqlizer Sqlizer) error {
	query, args, err := sqlizer.ToSql()
	if err != nil {
		return errors.WithStack(err)
	}
	_, err = execer.ExecContext(ctx, query, args...)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

type InsertDoc interface {
	Table() string
	Columns() []string
	Values() []any
}

func InsertOneDocReturningId(ctx context.Context, queryer QueryRowerContext, doc InsertDoc) (UUID, error) {
	sqlizer := sq.
		Insert(doc.Table()).
		Columns(doc.Columns()...).
		Values(doc.Values()...)
	return InsertValuesReturningId(ctx, queryer, sqlizer)
}

func InsertOneDoc(ctx context.Context, execer ExecerContext, doc InsertDoc) error {
	sqlizer := sq.
		Insert(doc.Table()).
		Columns(doc.Columns()...).
		Values(doc.Values()...)
	return InsertValues(ctx, execer, sqlizer)
}

func InsertManyDocs[D InsertDoc](ctx context.Context, execer ExecerContext, docs ...D) error {
	if len(docs) == 0 {
		return nil
	}
	sqlizer := sq.
		Insert(docs[0].Table()).
		Columns(docs[0].Columns()...)
	for _, doc := range docs {
		sqlizer = sqlizer.Values(doc.Values()...)
	}
	return InsertValues(ctx, execer, sqlizer)
}

func Exec(ctx context.Context, execer ExecerContext, sqlizer Sqlizer) (int64, error) {
	query, args, err := sqlizer.ToSql()
	if err != nil {
		return 0, errors.WithStack(err)
	}
	result, err := execer.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, errors.WithStack(err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, errors.WithStack(err)
	}
	return rowsAffected, nil
}

func FoundToNull[T any](val T, found bool, err error) (Null[T], error) {
	if err != nil {
		return EmptyNull[T](), err
	}
	if found {
		return FillNull(val), nil
	}
	return EmptyNull[T](), nil
}
