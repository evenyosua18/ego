package sqldb

import (
	"context"
	"github.com/evenyosua18/ego/sqldb"
)

type StubDbManager struct {
	// define return db or tx
	IsReturnDb bool

	// executor error
	ExecutorErr error

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

	// Begin tx
	BeginTxErr error

	// set db context error
	SetDBContextErr error

	StubExecutor
}

func (s *StubDbManager) GetExecutor(ctx context.Context) (sqldb.ISQLExecutor, error) {
	if s.ExecutorErr != nil {
		return nil, s.ExecutorErr
	}

	if !s.IsReturnDb {
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
			ExpectedQuery:     s.ExpectedQuery,
			ExpectedArgs:      s.ExpectedArgs,
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
			ExpectedQuery:     s.ExpectedQuery,
			ExpectedArgs:      s.ExpectedArgs,
		}, nil
	}
}

func (s *StubDbManager) BeginTx(ctx context.Context) (sqldb.ISqlTx, context.Context, error) {
	if s.BeginTxErr != nil {
		return nil, nil, s.BeginTxErr
	}

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
	}, ctx, nil
}

func (s *StubDbManager) SetDBContext(ctx context.Context) (context.Context, error) {
	if s.SetDBContextErr != nil {
		return ctx, s.SetDBContextErr
	}

	return ctx, nil
}
