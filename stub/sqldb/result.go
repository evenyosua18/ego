package sqldb

type StubSqlResult struct {
	Err   error
	Value int64
}

func (s *StubSqlResult) LastInsertId() (int64, error) {
	return s.Value, s.Err
}

func (s *StubSqlResult) RowsAffected() (int64, error) {
	return s.Value, s.Err
}
