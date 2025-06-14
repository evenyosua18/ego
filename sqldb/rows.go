package sqldb

import (
	"database/sql"
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

	return code.Wrap(err, code.DatabaseError)
}

func (s *SqlRows) ScanAll(dest any) error {
	err := dbscan.ScanAll(dest, s.rows)

	return code.Wrap(err, code.DatabaseError)
}

func (s *SqlRows) ScanOne(dest any) error {
	err := dbscan.ScanOne(dest, s.rows)

	return code.Wrap(err, code.DatabaseError)
}
