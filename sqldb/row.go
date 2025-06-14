package sqldb

import (
	"database/sql"
	"github.com/evenyosua18/ego/code"
	"github.com/pkg/errors"
)

type SqlRow struct {
	row *sql.Row
}

func NewSqlRow(row *sql.Row) ISqlRow {
	return &SqlRow{row: row}
}

func (s *SqlRow) Scan(dest ...interface{}) error {
	err := s.row.Scan(dest...)

	if errors.Is(err, sql.ErrNoRows) {
		return code.Wrap(err, code.NotFoundError)
	}

	return code.Wrap(err, code.DatabaseError)
}
