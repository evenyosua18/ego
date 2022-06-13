package postgresql

import "github.com/jinzhu/gorm"

type PostgresDriver struct {
	DB *gorm.DB
}

func NewPostgresDriver(db *gorm.DB) *PostgresDriver {
	return &PostgresDriver{DB: db}
}

func (d *PostgresDriver) Db() *gorm.DB {
	return d.DB
}
