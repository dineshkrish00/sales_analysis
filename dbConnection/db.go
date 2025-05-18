package dbConnection

import (
	"SalesReport/common"
	"SalesReport/global"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

// Structure to hold database connection details
type DatabaseType struct {
	Server   string
	Port     int
	User     string
	Password string
	Database string
	DBType   string
	DB       string
}

// structure to hold all db connection details used in this program
type AllUsedDatabases struct {
	Postgres DatabaseType
}

func Build_DB_Connection() error {
	log.Println("Build_DB_Connection (+)")
	var lErr error

	global.GConnection.DbPostgres, lErr = LocalDBConnection(Postgres)
	if lErr != nil {
		global.GConnection.DbPostgres.Close()
		common.ErrorMsg("Build_DB_Connection", lErr.Error())
		return lErr
	}

	log.Println("Build_DB_Connection (-)")
	return nil
}

func LocalDBConnection(DBtype string) (*sql.DB, error) {
	lDbDetails := new(AllUsedDatabases)
	lDbDetails.Init()
	lConnString := ""
	lLocalDBtype := ""

	var lDbCon *sql.DB
	var lErr error
	var lDBConnection DatabaseType

	if DBtype == lDbDetails.Postgres.DB {
		lDBConnection = lDbDetails.Postgres
		lLocalDBtype = lDbDetails.Postgres.DBType
	}

	// Build a Connection string for Postgres Connection Type
	if lLocalDBtype == "postgres" {
		lConnString = `user=` + lDBConnection.User + ` password=` + lDBConnection.Password + ` port=` + fmt.Sprintf("%v", lDBConnection.Port) + ` dbname=` + lDBConnection.Database + ` host=` + lDBConnection.Server + ` sslmode=disable`
	}

	//make a connection to db
	if lLocalDBtype != "" {

		lDbCon, lErr = sql.Open(lLocalDBtype, lConnString)
		if lErr != nil {
			log.Println("Open connection failed:", lErr.Error())
		} else {

			Connections := common.ReadTomlConfig("./toml/dbconfig.toml")

			global.GConfig.Db_Max_OpenConn = fmt.Sprintf("%v", Connections.(map[string]interface{})["DB_Max_Open_Connection"])
			global.GConfig.Db_Max_IdleConn = fmt.Sprintf("%v", Connections.(map[string]interface{})["DB_Max_Idle_Connection"])
			global.GConfig.Db_ConnMax_IdleTime = fmt.Sprintf("%v", Connections.(map[string]interface{})["DB_Max_Idle_Time"])

			//Converting String to Integer
			Db_Max_OpenConn, _ := strconv.Atoi(global.GConfig.Db_Max_OpenConn)
			Db_Max_IdleConn, _ := strconv.Atoi(global.GConfig.Db_Max_IdleConn)
			Db_ConnMax_IdleTime, _ := strconv.Atoi(global.GConfig.Db_ConnMax_IdleTime)

			lDbCon.SetMaxOpenConns(Db_Max_OpenConn)
			lDbCon.SetMaxIdleConns(Db_Max_IdleConn)
			lDbCon.SetConnMaxIdleTime(time.Second * time.Duration(Db_ConnMax_IdleTime))
		}
	} else {
		return lDbCon, fmt.Errorf(" Invalid DB Details Credentials")
	}
	return lDbCon, lErr
}
