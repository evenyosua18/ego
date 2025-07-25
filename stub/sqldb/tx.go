package sqldb

import (
	"fmt"
	"github.com/evenyosua18/ego/sqldb"
)

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

	// expected query
	ExpectedQuery string
	ExpectedArgs  []any
}

func (s *StubSqlTx) QueryRow(query string, args ...any) sqldb.ISqlRow {
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

func (s *StubSqlTx) Query(query string, args ...any) (sqldb.ISqlRows, error) {
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

func (s *StubSqlTx) Exec(query string, args ...any) (sqldb.ISqlResult, error) {
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
