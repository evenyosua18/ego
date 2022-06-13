package database

import (
	"github.com/jinzhu/gorm"
)

type DbDriver interface {
	Db() *gorm.DB
}

//type DML interface {
//}
//
//type Pagination interface {
//}
