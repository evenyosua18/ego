package sqldb

import (
	"database/sql"
	"errors"
	"github.com/evenyosua18/ego/code"
	"github.com/georgysavva/scany/v2/dbscan"
)

type SqlRows struct {
	rows *sql.Rows
}

func NewSqlRows(rows *sql.Rows) ISqlRows {
	return &SqlRows{rows: rows}
}

func (s *SqlRows) Close() error {
	return s.rows.Close()
}

func (s *SqlRows) Next() bool {
	return s.rows.Next()
}

func (s *SqlRows) MapResult(dest ...any) error {
	err := s.rows.Scan(dest...)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return code.Wrap(err, code.NotFoundError)
		}

		return code.Wrap(err, code.DatabaseError)
	}

	return nil
}

func (s *SqlRows) ScanAll(dest any) error {
	err := dbscan.ScanAll(dest, s.rows)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return code.Wrap(err, code.NotFoundError)
		}

		return code.Wrap(err, code.DatabaseError)
	}

	return nil
}

func (s *SqlRows) ScanOne(dest any) error {
	err := dbscan.ScanOne(dest, s.rows)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return code.Wrap(err, code.NotFoundError)
		}

		return code.Wrap(err, code.DatabaseError)
	}

	return nil
}
