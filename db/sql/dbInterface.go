package sql

import "gorm.io/gorm"

type DbDriver interface {
	DB() *gorm.DB
}
