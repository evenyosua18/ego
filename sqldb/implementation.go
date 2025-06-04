package sqldb

import (
	"context"
	"fmt"
)

type (
	dbKey struct{}
	txKey struct{}
)

var (
	errNoConnection       = fmt.Errorf("no connection found")
	errNoDB               = fmt.Errorf("no database found")
	errDBAlreadyConnected = fmt.Errorf("database already connected")
)

func GetExecutor(ctx context.Context) (ISQLExecutor, error) {
	if tx, ok := ctx.Value(txKey{}).(ISqlTx); ok {
		return tx, nil
	}

	if db, ok := ctx.Value(dbKey{}).(ISqlDB); ok {
		return db, nil
	}

	return nil, errNoConnection
}

func BeginTx(ctx context.Context) (ISqlTx, context.Context, error) {
	// get sql db
	sqlDb, err := GetDB()

	if err != nil {
		return nil, nil, err
	}

	// begin tx
	tx, err := sqlDb.Begin()

	if err != nil {
		return nil, nil, err
	}

	// set to context
	newCtx := context.WithValue(ctx, txKey{}, tx)

	return tx, newCtx, nil
}

func SetDBContext(ctx context.Context) (context.Context, error) {
	// prevent empty db
	sqlDb, err := GetDB()

	if err != nil {
		return nil, err
	}

	newCtx := context.WithValue(ctx, dbKey{}, sqlDb)

	return newCtx, nil
}
