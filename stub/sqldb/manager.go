package sqldb

import (
	"context"
	"github.com/evenyosua18/ego/sqldb"
)

type StubDbManager struct {
	// manager
	IsReturnTx bool
	ManagerErr error

	// commit
	CommitErr error

	// rollback
	RollbackErr error

	// end tx
	EndTxErr error

	// Ping
	PingErr error

	// Begin
	BeginErr error

	StubExecutor
}

func (s *StubDbManager) GetExecutor(ctx context.Context) (sqldb.ISQLExecutor, error) {
	if s.ManagerErr != nil {
		return nil, s.ManagerErr
	}

	if s.IsReturnTx {
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
		}, nil
	} else {
		return &StubSqlDb{
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
			PingErr:           s.PingErr,
			BeginErr:          s.BeginErr,
			CommitErr:         s.CommitErr,
			RollbackErr:       s.RollbackErr,
			EndTxErr:          s.EndTxErr,
		}, nil
	}
}
