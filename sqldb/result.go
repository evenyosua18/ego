package sqldb

import (
	"database/sql"
	"github.com/evenyosua18/ego/code"
)

type SqlResult struct {
	result sql.Result
}

func NewSqlResult(result sql.Result) ISqlResult {
	return &SqlResult{result: result}
}

func (s *SqlResult) LastInsertId() (int64, error) {
	lastInsertedId, err := s.result.LastInsertId()

	return lastInsertedId, code.Wrap(err, code.DatabaseError)
}

func (s *SqlResult) RowsAffected() (int64, error) {
	rowsAffected, err := s.result.RowsAffected()

	return rowsAffected, code.Wrap(err, code.DatabaseError)
}
