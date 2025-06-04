package sqldb

import "database/sql"

type SqlResult struct {
	result sql.Result
}

func NewSqlResult(result sql.Result) ISqlResult {
	return &SqlResult{result: result}
}

func (s *SqlResult) LastInsertId() (int64, error) {
	return s.result.LastInsertId()
}

func (s *SqlResult) RowsAffected() (int64, error) {
	return s.result.RowsAffected()
}
