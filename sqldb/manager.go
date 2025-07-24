package sqldb

import (
	"context"
	"fmt"
	"github.com/evenyosua18/ego/code"
)

type (
	dbKey struct{}
	txKey struct{}
)

var (
	dbManager IDbManager = &DbManager{}

	errNoConnection       = code.Get(code.DatabaseError).SetErrorMessage("no connection found")
	errNoDB               = code.Get(code.DatabaseError).SetErrorMessage("no database found")
	errDBAlreadyConnected = fmt.Errorf("database already connected")
)

type DbManager struct{}

func (s *DbManager) GetExecutor(ctx context.Context) (ISQLExecutor, error) {
	if tx, ok := ctx.Value(txKey{}).(ISqlTx); ok {
		return tx, nil
	}

	if db, ok := ctx.Value(dbKey{}).(ISqlDB); ok {
		return db, nil
	}

	return nil, errNoConnection
}

func (s *DbManager) BeginTx(ctx context.Context) (ISqlTx, context.Context, error) {
	// get sql db
	sqlDb, err := GetDB()

	if err != nil {
		return nil, nil, err
	}

	// begin tx
	tx, err := sqlDb.Begin()

	if err != nil {
		return nil, nil, code.Wrap(err, code.DatabaseError)
	}

	// set to context
	newCtx := context.WithValue(ctx, txKey{}, tx)

	return tx, newCtx, nil
}

func (s *DbManager) SetDBContext(ctx context.Context) (context.Context, error) {
	// prevent empty db
	sqlDb, err := GetDB()

	if err != nil {
		return nil, err
	}

	newCtx := context.WithValue(ctx, dbKey{}, sqlDb)

	return newCtx, nil
}

func GetDbManager() IDbManager {
	return dbManager
}
