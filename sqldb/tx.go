package sqldb

import (
	"database/sql"
	"sync"
)

type SqlTx struct {
	tx   *sql.Tx
	once sync.Once
}

func NewSqlTx(tx *sql.Tx) ISqlTx {
	return &SqlTx{tx: tx}
}

func (s *SqlTx) QueryRow(query string, args ...any) ISqlRow {
	row := s.tx.QueryRow(query, args...)
	return NewSqlRow(row)
}

func (s *SqlTx) Query(query string, args ...any) (ISqlRows, error) {
	rows, err := s.tx.Query(query, args...)

	if err != nil {
		return nil, err
	}

	return NewSqlRows(rows), nil
}

func (s *SqlTx) Exec(query string, args ...any) (ISqlResult, error) {
	rows, err := s.tx.Exec(query, args...)

	if err != nil {
		return nil, err
	}

	return NewSqlResult(rows), nil
}

func (s *SqlTx) Rollback() error {
	return s.tx.Rollback()
}

func (s *SqlTx) Commit() error {
	return s.tx.Commit()
}

func (s *SqlTx) EndTx(errPtr *error) {
	defer s.once.Do(func() {
		if r := recover(); r != nil {
			_ = s.tx.Rollback()
			panic(r)
		}

		if *errPtr != nil {
			_ = s.tx.Rollback()
			return
		}

		_ = s.tx.Commit()
	})
}
