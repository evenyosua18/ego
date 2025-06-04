package sqldb

import (
	"database/sql"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type SqlDB struct {
	db *sql.DB
}

type Config struct {
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

var (
	mu sync.RWMutex
	db ISqlDB
)

// Connect creates a new SQL DB connection with the provided configuration.
// it does NOT set the global db. Use SetDB() if global access is needed.
func Connect(driver, uri string, config *Config) (ISqlDB, error) {
	if db != nil {
		return nil, errDBAlreadyConnected
	}

	// initialize db here
	dbConnection, err := sql.Open(driver, uri)
	if err != nil {
		return nil, err
	}

	// config
	if config != nil {
		if config.MaxOpenConns <= 0 {
			dbConnection.SetMaxOpenConns(config.MaxOpenConns)
		}

		if config.MaxIdleConns <= 0 {
			dbConnection.SetMaxIdleConns(config.MaxIdleConns)
		}

		if config.ConnMaxLifetime != 0 {
			dbConnection.SetConnMaxLifetime(config.ConnMaxLifetime)
		}

		if config.ConnMaxIdleTime != 0 {
			dbConnection.SetConnMaxIdleTime(config.ConnMaxIdleTime)
		}
	}

	// ping
	if err = dbConnection.Ping(); err != nil {
		return nil, err
	}

	return &SqlDB{db: dbConnection}, nil
}

func SetDB(sqlDB ISqlDB) {
	mu.Lock()
	defer mu.Unlock()

	db = sqlDB
}

func GetDB() (ISqlDB, error) {
	mu.RLock()
	defer mu.RUnlock()

	if db == nil {
		return nil, errNoDB
	}

	return db, nil
}

func ResetDB() {
	mu.Lock()
	defer mu.Unlock()

	db = nil
}

func (s *SqlDB) QueryRow(query string, args ...any) ISqlRow {
	row := s.db.QueryRow(query, args...)
	return NewSqlRow(row)
}

func (s *SqlDB) Query(query string, args ...any) (ISqlRows, error) {
	rows, err := s.db.Query(query, args...)

	if err != nil {
		return nil, err
	}

	return NewSqlRows(rows), nil
}

func (s *SqlDB) Exec(query string, args ...any) (ISqlResult, error) {
	rows, err := s.db.Exec(query, args...)

	if err != nil {
		return nil, err
	}

	return NewSqlResult(rows), nil
}

func (s *SqlDB) Begin() (ISqlTx, error) {
	tx, err := s.db.Begin()

	if err != nil {
		return nil, err
	}

	return NewSqlTx(tx), nil
}

func (s *SqlDB) Ping() error {
	return s.db.Ping()
}
