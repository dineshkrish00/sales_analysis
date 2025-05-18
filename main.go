package main

import (
	"SalesReport/common"
	fileoperation "SalesReport/dataintegration/csv_file_operation"
	fetchproducts "SalesReport/dataintegration/fetch_products"
	fetchrevenue "SalesReport/dataintegration/fetch_revenue"
	"strconv"

	"SalesReport/dbConnection"
	"net/http"

	"SalesReport/global"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	log.Println("Server Execution Started (+)")

	lMethod := "main"

	//create log file
	f, lErr := os.OpenFile("./log/logfile"+time.Now().Format("02012006.15.04.05.000000000")+".txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if lErr != nil {
		log.Fatalf("error opening file: %v", lErr)
	}
	defer f.Close()

	log.SetOutput(f)

	lErr = dbConnection.Build_DB_Connection()
	if lErr != nil {
		log.Println("Error in Build_DB_Connection", lErr.Error())
		common.ErrorMsg(lMethod, lErr.Error())
	} else {

		defer global.GConnection.DbPostgres.Close()

		lStatus := common.ReadTomlConfig("./toml/config.toml")
		lTrigger := fmt.Sprintf("%v", lStatus.(map[string]interface{})["TriggerStatus"])

		if lTrigger == "Y" {
			go CsvDataRefresh()
		}

		http.HandleFunc("/csvfileupload", fileoperation.ReadandInsertData)
		// http.HandleFunc("/upload", fileoperation.Testing)
		http.HandleFunc("/fetchRevenue", fetchrevenue.GetCategoryWiseRevenue)
		http.HandleFunc("/fetchProductDetails", fetchproducts.GetCategoryWiseProduct)

		http.ListenAndServe(":29001", nil)
		log.Println("Server Execution Started (-)")
	}
}

func CsvDataRefresh() {
	lastRunDate := ""

	for {
		now := time.Now()
		currentDate := now.Format("2006-01-02") // yyyy-mm-dd

		lStatus := common.ReadTomlConfig("./toml/config.toml")
		lHour := fmt.Sprintf("%v", lStatus.(map[string]interface{})["ScheduleHour"])
		lMinute := fmt.Sprintf("%v", lStatus.(map[string]interface{})["ScheduleMinute"])

		lSHour, _ := strconv.Atoi(lHour)

		lSMinute, _ := strconv.Atoi(lMinute)

		// Check if time is 5:30 AM and
		if now.Hour() == lSHour && now.Minute() == lSMinute && lastRunDate != currentDate {
			log.Println("Running DataRefreshMechanism at 5:30 AM...")
			DataRefreshMechanism()
			lastRunDate = currentDate
		}

		time.Sleep(1 * time.Minute)
	}
}

func DataRefreshMechanism() {
	log.Println("DataRefreshMechanism (+)")
	lErr := fileoperation.ExtractandInstert()
	if lErr != nil {
		common.ErrorMsg("ReadandInsertData", lErr.Error())
	}
	log.Println("DataRefreshMechanism (-)")
}
