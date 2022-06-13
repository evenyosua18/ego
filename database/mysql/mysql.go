package mysql

import (
	"github.com/jinzhu/gorm"
)

type MysqlDriver struct {
	DB *gorm.DB
}

func NewMysqlDriver(db *gorm.DB) *MysqlDriver {
	return &MysqlDriver{DB: db}
}

func (d *MysqlDriver) Db() *gorm.DB {
	return d.DB
}
