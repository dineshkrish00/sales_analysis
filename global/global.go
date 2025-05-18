package global

import (
	"database/sql"
)

type DbConnection struct {
	DbPostgres *sql.DB
}

var GConnection DbConnection

type ConfigToml struct {
	Db_Max_OpenConn     string
	Db_Max_IdleConn     string
	Db_ConnMax_IdleTime string
}

var GConfig ConfigToml
