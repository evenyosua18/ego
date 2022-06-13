package database

import (
	"ego/database/mysql"
	"ego/database/postgresql"
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Database struct {
	Username string `yaml:"username"` //username of account for login into database
	Password string `yaml:"password"` //password of account for login into database
	Port     string `yaml:"port"`     //port of database that has been stored
	Address  string `yaml:"address"`  //ip address of database that has been stored
	Database string `yaml:"database"` //name of database
	Adapter  string `yaml:"adapter"`  //name of adapter that use on the struct
}

type DBConfig struct {
}

const (
	DbMysqlAdapter     = `mysql`
	DbPsqlAdapter      = `postgresql`
	DbSqlServerAdapter = `mssql`

	DriverNotRecognize = `not recognize any db driver`
)

func NewDatabaseInstance(database Database) (DbDriver, error) {
	switch database.Adapter {
	case DbMysqlAdapter:
		return newMysqlDriver(database), nil
	case DbPsqlAdapter:
		return newPsqlDriver(database), nil
		break
	case DbSqlServerAdapter:
		break
	default:
		return nil, errors.New(DriverNotRecognize)
	}

	return nil, errors.New(DriverNotRecognize)
}

//New Mysql Driver charset, parseTime, local, etc can be more manageable with configuration
func newMysqlDriver(database Database) DbDriver {
	dsn := database.Username + ":" + database.Password + "@tcp(" + database.Address + ":" + database.Port + ")/" + database.Username + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(database.Adapter, dsn)
	if err != nil {
		return nil
	}

	return mysql.NewMysqlDriver(db)
}

func newPsqlDriver(database Database) DbDriver {
	dsn := "host=" + database.Address + " port=" + database.Port + " user=" + database.Username + " dbname=" + database.Database + " password=" + database.Password
	db, err := gorm.Open(database.Adapter, dsn)

	if err != nil {
		return nil
	}

	return postgresql.NewPostgresDriver(db)
}
