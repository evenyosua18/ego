package sqldb

import (
	"database/sql"
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
	return s.rows.Scan(dest...)
}

func (s *SqlRows) ScanAll(dest any) error {
	return dbscan.ScanAll(dest, s.rows)
}

func (s *SqlRows) ScanOne(dest any) error {
	return dbscan.ScanOne(dest, s.rows)
}
