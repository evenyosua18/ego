package sql

import "gorm.io/gorm"

type DatabaseConfig struct {
	Username string `yaml:"username"` //username of account for login into database
	Password string `yaml:"password"` //password of account for login into database
	Port     string `yaml:"port"`     //port of database that has been stored
	Address  string `yaml:"address"`  //ip address of database that has been stored
	Database string `yaml:"database"` //name of database
	Adapter  string `yaml:"adapter"`  //name of adapter that use on the struct

	SSLMode string `yaml:"sslmode"` //ssl mode for db connection
}

type DbModel struct {
	DB *gorm.DB
}

type (
	PostgresSSLMode string
)

const (
	DbMysqlAdapter     = `mysql`
	DbPsqlAdapter      = `postgres`
	DbSqlServerAdapter = `mssql`

	DriverNotRecognize = `not recognize any db driver`

	SSLModeDisable    PostgresSSLMode = `disable`
	SSLModeAllow      PostgresSSLMode = `allow`
	SSLModePrefer     PostgresSSLMode = `prefer`
	SSLModeRequire    PostgresSSLMode = `require`
	SSLModeVerifyCA   PostgresSSLMode = `verify-ca`
	SSLModeVerifyFull PostgresSSLMode = `verify-full`
)
