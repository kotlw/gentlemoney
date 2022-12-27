package sqlite

import (
	"database/sql"
	"fmt"

	"gentlemoney/internal/model"
)

// executor is a wrapper for sql.DB.Exec() and sql.DB.Query().
type executor[T model.Any] struct {
	db *sql.DB
}

// insert executes insert query with given arguments.
func (e *executor[_]) insert(query string, args ...any) (int64, error) {
	res, err := e.db.Exec(query, args...)
	if err != nil {
        return -1, fmt.Errorf("e.db.Exec: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
        return -1, fmt.Errorf("res.LastInsertId: %w", err)
	}

	return id, nil
}

// update executes update and delete queries with given arguments.
func (e *executor[_]) update(query string, args ...any) error {
	res, err := e.db.Exec(query, args...)
	if err != nil {
        return fmt.Errorf("e.db.Exec: %w", err)
	}

	rowsAfected, err := res.RowsAffected()
	if err != nil {
        return fmt.Errorf("res.RowsAffected: %w", err)
	}
	if rowsAfected != 1 {
		return fmt.Errorf("total affected rows %d while expected 1", rowsAfected)
	}

	return nil
}

// getAll returns all rows from persistant storage, it requires dest func which should return new
// object of certain type, and addreses of its fields to Scan. Order of addreses should match with
// order of coresponding columns in query.
func (e *executor[T]) getAll(query string, dest func() (*T, []any)) ([]*T, error) {
	rows, err := e.db.Query(query)
	if err != nil {
        return nil, fmt.Errorf("stmt.Query: %w", err)
	}
	defer func() {
		if err = rows.Close(); err != nil {
            err = fmt.Errorf("defer rows.Close: %w", err)
		}
	}()

	res := make([]*T, 0, 20)
	for rows.Next() {
		t, addrs := dest()
		if err = rows.Scan(addrs...); err != nil {
            return nil, fmt.Errorf("rows.Scan: %w", err)
		}
		res = append(res, t)
	}

	return res, err
}
