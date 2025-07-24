package sqldb

import "github.com/evenyosua18/ego/sqldb"

type StubSqlTx struct {
	ExecError error

	// result
	ResultValue int64
	ResultError error

	// query row
	QueryRowValues []any
	QueryRowErr    error

	// query
	QueryValues      []any
	QueryDestination any

	QueryErr          error
	QueryCloseErr     error
	QueryMapResultErr error
	QueryScanAllErr   error
	QueryScanOneErr   error

	QueryRowIndex int

	// commit
	CommitErr error

	// rollback
	RollbackErr error

	// end tx
	EndTxErr error
}

func (s *StubSqlTx) QueryRow(query string, args ...any) sqldb.ISqlRow {
	return &StubSqlRow{Values: s.QueryRowValues, Err: s.QueryRowErr}
}

func (s *StubSqlTx) Query(query string, args ...any) (sqldb.ISqlRows, error) {
	return &StubSqlRows{
		Values:       s.QueryValues,
		Destination:  s.QueryDestination,
		CloseErr:     s.QueryCloseErr,
		MapResultErr: s.QueryMapResultErr,
		ScanAllErr:   s.QueryScanAllErr,
		ScanOneErr:   s.QueryScanOneErr,
		RowIndex:     s.QueryRowIndex,
	}, s.QueryErr
}

func (s *StubSqlTx) Exec(query string, args ...any) (sqldb.ISqlResult, error) {
	return &StubSqlResult{Value: s.ResultValue, Err: s.ResultError}, s.ExecError
}

func (s *StubSqlTx) Rollback() error {
	return s.RollbackErr
}

func (s *StubSqlTx) Commit() error {
	return s.CommitErr
}

func (s *StubSqlTx) EndTx(ptrErr *error) {
	if s.EndTxErr != nil {
		*ptrErr = s.EndTxErr
	}

	return
}
