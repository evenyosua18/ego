package postgres

import (
	"github.com/evenyosua18/ego/db/sql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB(dbConfig sql.DatabaseConfig) (*sql.DbModel, error) {
	//open connection
	additionalConfig := " "

	if dbConfig.SSLMode != "" {
		additionalConfig += "sslmode=" + string(dbConfig.SSLMode) + " "
	} else {
		additionalConfig += "sslmode=" + string(sql.SSLModeDisable) + " "
	}

	dsn := "host=" + dbConfig.Address + " port=" + dbConfig.Port + " user=" + dbConfig.Username + " dbname=" + dbConfig.Database + " password=" + dbConfig.Password + additionalConfig

	db, err := gorm.Open(postgres.Open(dsn))

	if err != nil {
		return nil, err
	}

	return &sql.DbModel{
		DB: db,
	}, nil
}
