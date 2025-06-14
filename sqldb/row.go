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

func (s *SqlRow) Scan(dest ...any) error {
	err := s.row.Scan(dest...)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return code.Wrap(err, code.NotFoundError)
		}

		return code.Wrap(err, code.DatabaseError)
	}

	return nil
}
