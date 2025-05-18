package fileoperation

import (
	"SalesReport/common"
	"SalesReport/model"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

/*
Method: Purpose To get the CSV file from api end poin or path triggering it will process and insert
into respective table
*/
func ReadandInsertData(w http.ResponseWriter, req *http.Request) {
	log.Println("ReadandInsertData (+)")
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "POST,OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if req.Method == http.MethodPost {

		var lRespRec model.ResponseStruct

		lRespRec.Status = common.SuccessCode

		// Check if the request has a file it will process with tath
		if req.ContentLength > 0 {
			lRecords, lErr := ExtractDataRequest(req)
			if lErr != nil {
				lRespRec.Status = common.ErrorCode
				common.ErrorMsg("ExtractDataRequest", lErr.Error())
			} else {
				lErr = FileReadAndInsert(lRecords)
				if lErr != nil {
					lRespRec.Status = common.ErrorCode
					common.ErrorMsg("FileReadAndInsert", lErr.Error())
				} else {
					lRespRec.Msg = "File Uploaded Succesfully"
				}
			}
		} else {
			// No file in the request it execute the direct file from the path
			lErr := ExtractandInstert()
			if lErr != nil {
				lRespRec.Status = common.ErrorCode
				common.ErrorMsg("ExtractandInstert", lErr.Error())
			} else {
				lRespRec.Msg = "File Uploaded Succesfully"
			}
		}

		// After processing, check for any final errors
		lData, lErr := json.Marshal(lRespRec)
		if lErr != nil {
			common.ErrorMsg("ReadandInsertData", lErr.Error())
			fmt.Fprint(w, lErr.Error())
		} else {
			fmt.Fprint(w, string(lData))
		}
	}
	log.Println("ReadandInsertData (-)")
}
