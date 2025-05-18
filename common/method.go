package common

import (
	"log"
	"runtime"
	"strconv"
	"time"

	"github.com/BurntSushi/toml"
)

// --------------------------------------------------------------------
// function reads the constants from the config.toml file
// --------------------------------------------------------------------

func ReadTomlConfig(filename string) interface{} {
	var folder interface{}
	if _, lErr := toml.DecodeFile(filename, &folder); lErr != nil {
		ErrorMsg("ReadTomlConfig", lErr.Error())

	}
	return folder
}

// Errmsg Log with line no and Method name will capture
func ErrorMsg(MethodName, Err_Msg string) {
	_, _, number, ok := runtime.Caller(1)
	if !ok {
		log.Println("Error  line number", number)
	}
	log.Println("Err Msg :" + Err_Msg + " (Line No " + strconv.Itoa(number) + " and  Method Name is " + MethodName + ")")
}

// Converting Date into Standard
func GetDate(pdateStr string) time.Time {

	lDate, lErr := time.Parse("2006-01-02", pdateStr)
	if lErr != nil {
		ErrorMsg("GetDate", lErr.Error())
	}

	return lDate
}
