package sqldb

import "database/sql"

type SqlRow struct {
	row *sql.Row
}

func NewSqlRow(row *sql.Row) ISqlRow {
	return &SqlRow{row: row}
}

func (s *SqlRow) Scan(dest ...interface{}) error {
	return s.row.Scan(dest...)
}
