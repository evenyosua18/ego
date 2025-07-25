package sqldb

import (
	"fmt"
	"github.com/evenyosua18/ego/sqldb"
)

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

	// expected parameter
	ExpectedQuery string
	ExpectedArgs  []any
}

func (s *StubSqlDb) QueryRow(query string, args ...any) sqldb.ISqlRow {
	if query != s.ExpectedQuery {
		return &StubSqlRow{Values: nil, Err: fmt.Errorf(`unexpected query %s, want %s`, query, s.ExpectedQuery)}
	}

	if len(args) != len(s.ExpectedArgs) {
		return &StubSqlRow{Values: nil, Err: fmt.Errorf(`unexpected length of args %d, want %d`, len(args), len(s.ExpectedArgs))}
	}

	for i := 0; i < len(args); i++ {
		if args[i] != s.ExpectedArgs[i] {
			return &StubSqlRow{Values: nil, Err: fmt.Errorf(`unexpected value of arg %v, want %v`, args[i], s.ExpectedArgs[i])}
		}
	}

	return &StubSqlRow{Values: s.QueryRowValues, Err: s.QueryRowErr}
}

func (s *StubSqlDb) Query(query string, args ...any) (sqldb.ISqlRows, error) {
	if query != s.ExpectedQuery {
		return &StubSqlRows{}, fmt.Errorf(`unexpected query %s, want %s`, query, s.ExpectedQuery)
	}

	if len(args) != len(s.ExpectedArgs) {
		return &StubSqlRows{}, fmt.Errorf(`unexpected length of args %d, want %d`, len(args), len(s.ExpectedArgs))
	}

	for i := 0; i < len(args); i++ {
		if args[i] != s.ExpectedArgs[i] {
			return &StubSqlRows{}, fmt.Errorf(`unexpected value of arg %v, want %v`, args[i], s.ExpectedArgs[i])
		}
	}

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
	if query != s.ExpectedQuery {
		return &StubSqlResult{}, fmt.Errorf(`unexpected query %s, want %s`, query, s.ExpectedQuery)
	}

	if len(args) != len(s.ExpectedArgs) {
		return &StubSqlResult{}, fmt.Errorf(`unexpected length of args %d, want %d`, len(args), len(s.ExpectedArgs))
	}

	for i := 0; i < len(args); i++ {
		if args[i] != s.ExpectedArgs[i] {
			return &StubSqlResult{}, fmt.Errorf(`unexpected value of arg %v, want %v`, args[i], s.ExpectedArgs[i])
		}
	}

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
