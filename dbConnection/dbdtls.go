package dbConnection

import (
	"SalesReport/common"
	"fmt"
	"strconv"
)

const (
	Postgres = "POSTGRES"
)

// Initializing DB Details
func (d *AllUsedDatabases) Init() {
	dbconfig := common.ReadTomlConfig("./toml/dbconfig.toml")

	//setting IPO db connection details
	d.Postgres.Server = fmt.Sprintf("%v", dbconfig.(map[string]interface{})["PostgresDBServer"])
	d.Postgres.Port, _ = strconv.Atoi(fmt.Sprintf("%v", dbconfig.(map[string]interface{})["PostgresDBPort"]))
	d.Postgres.User = fmt.Sprintf("%v", dbconfig.(map[string]interface{})["PostgresDBUser"])
	d.Postgres.Password = fmt.Sprintf("%v", dbconfig.(map[string]interface{})["PostgresDBPassword"])
	d.Postgres.Database = fmt.Sprintf("%v", dbconfig.(map[string]interface{})["PostgresDBDatabase"])
	d.Postgres.DBType = fmt.Sprintf("%v", dbconfig.(map[string]interface{})["PostgresDBType"])
	d.Postgres.DB = Postgres
}
