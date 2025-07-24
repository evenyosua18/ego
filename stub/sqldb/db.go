package sqldb

import "github.com/evenyosua18/ego/sqldb"

type StubSqlDb struct {
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

	// Ping
	PingErr error

	// Begin
	BeginErr error

	// commit
	CommitErr error

	// rollback
	RollbackErr error

	// end tx
	EndTxErr error
}

func (s *StubSqlDb) QueryRow(query string, args ...any) sqldb.ISqlRow {
	return &StubSqlRow{Values: s.QueryRowValues, Err: s.QueryRowErr}
}

func (s *StubSqlDb) Query(query string, args ...any) (sqldb.ISqlRows, error) {
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

func (s *StubSqlDb) Exec(query string, args ...any) (sqldb.ISqlResult, error) {
	return &StubSqlResult{Value: s.ResultValue, Err: s.ResultError}, s.ExecError
}

func (s *StubSqlDb) Begin() (sqldb.ISqlTx, error) {
	return &StubSqlTx{
		ExecError:         s.ExecError,
		ResultValue:       s.ResultValue,
		ResultError:       s.ResultError,
		QueryRowValues:    s.QueryRowValues,
		QueryRowErr:       s.QueryRowErr,
		QueryValues:       s.QueryValues,
		QueryDestination:  s.QueryDestination,
		QueryErr:          s.QueryErr,
		QueryCloseErr:     s.QueryCloseErr,
		QueryMapResultErr: s.QueryMapResultErr,
		QueryScanAllErr:   s.QueryScanAllErr,
		QueryScanOneErr:   s.QueryScanOneErr,
		QueryRowIndex:     s.QueryRowIndex,
		CommitErr:         s.CommitErr,
		RollbackErr:       s.RollbackErr,
		EndTxErr:          s.EndTxErr,
	}, s.BeginErr
}

func (s *StubSqlDb) Ping() error {
	return s.PingErr
}
