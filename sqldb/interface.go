package sqldb

type ISqlDB interface {
	ISQLExecutor
	Begin() (ISqlTx, error)
	Ping() error
}

type ISqlTx interface {
	ISQLExecutor
	Rollback() error
	Commit() error
	EndTx(ptrErr *error)
}

type ISqlRow interface {
	Scan(dest ...any) error
}

type ISqlRows interface {
	Close() error
	Next() bool
	MapResult(dest ...any) error
	ScanAll(dest any) error
	ScanOne(dest any) error
}

type ISqlResult interface {
	LastInsertId() (int64, error)
	RowsAffected() (int64, error)
}

type ISQLExecutor interface {
	QueryRow(query string, args ...any) ISqlRow
	Query(query string, args ...any) (ISqlRows, error)
	Exec(query string, args ...any) (ISqlResult, error)
}
