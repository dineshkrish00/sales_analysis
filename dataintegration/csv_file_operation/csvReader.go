package fileoperation

import (
	"SalesReport/common"
	"SalesReport/model"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

//Read the csv file from the form data

func ExtractDataRequest(r *http.Request) ([][]string, error) {
	log.Println("CSVRequestCheck (+)")

	lFile, _, lErr := r.FormFile("csvfile")
	if lErr == nil {
		defer lFile.Close()
		log.Println("Reading CSV from uploaded file ")
		return ReadCSVFile(lFile)
	}
	return ReadCSVFile(lFile)
}

// Get Records from api and Extract then process on this method
func FileReadAndInsert(lRecords [][]string) error {
	log.Println("FileReadAndInsert (+)")

	//Split and Insert into each Seperate Table
	lErr := SplitAndInsert(lRecords)
	if lErr != nil {
		common.ErrorMsg("FileReadAndInsert", lErr.Error())
		return lErr
	}
	log.Println("FileReadAndInsert (-)")
	return nil
}

func ExtractandInstert() error {
	log.Println("ExtractandInstert (+)")

	// Read the CSV file data and store it in a slice
	// of string slices
	lRecords, lErr := FileReader()
	if lErr != nil {
		common.ErrorMsg("ExtractandInstert", lErr.Error())
		return lErr
	} else {
		//Split and Insert into each Seperate Table
		lErr = SplitAndInsert(lRecords)
		if lErr != nil {
			common.ErrorMsg("ExtractandInstert", lErr.Error())
			return lErr
		}
	}
	log.Println("ExtractandInstert (-)")
	return nil
}

// /Read the file path and returns the [][]string
func FileReader() ([][]string, error) {

	log.Println("FileReader (+)")
	lpathConfig := common.ReadTomlConfig("./toml/fileconfig.toml")
	lfilepath := fmt.Sprintf("%v", lpathConfig.(map[string]interface{})["FilePath"])
	file, lErr := os.Open(lfilepath)
	if lErr != nil {
		common.ErrorMsg("Open", lErr.Error())
		return nil, lErr
	}
	defer file.Close()

	log.Println("FileReader (-)")
	return ReadCSVFile(file)
}

// ReadCSVFile reads and parses a CSV file from the given reader (can be file or uploaded file)
func ReadCSVFile(reader io.Reader) ([][]string, error) {
	// Create a CSV reader from the provided reader (file, etc.)
	lCsvReader := csv.NewReader(reader)
	lRecords, lErr := lCsvReader.ReadAll()
	if lErr != nil {
		common.ErrorMsg("ReadCSVFile", lErr.Error())
		return nil, lErr
	}
	return lRecords, nil
}

func SplitAndInsert(pInputRecords [][]string) error {
	log.Println("SplitAndInsert (+)")

	// Iterate through the records and insert them into the database

	for idx, lRow := range pInputRecords {
		if idx == 0 {
			continue
		}

		// Customer Details Check if record exist if exist not only insert the record into customers table
		lCustomerID, lErr := Check_Record_If_Exists("customers", "customer_id", lRow[2])
		if lErr != nil {
			common.ErrorMsg("Check_If_RecordExists", lErr.Error())
			return lErr
		}
		if lCustomerID == 0 {
			lCustomerID, lErr = Insert_Customer_Data(model.CustomerDataStruct{
				Customer_id:      lRow[2],
				Customer_name:    lRow[12],
				Customer_email:   lRow[13],
				Customer_address: lRow[14],
			})
			if lErr != nil {
				common.ErrorMsg("Insert_Customer_Data", lErr.Error())
				return lErr
			}
		}

		// Product Details Check if record exist if exist not only insert the record into products table
		lProductID, lErr := Check_Record_If_Exists("products", "product_id", lRow[1])
		if lErr != nil {
			common.ErrorMsg("Check_If_RecordExists", lErr.Error())
			return lErr
		}
		if lProductID == 0 {
			lProductID, lErr = Insert_Product_Data(model.ProductsStruct{
				Product_id:   lRow[1],
				Product_Name: lRow[3],
				Category:     lRow[4],
			})
			if lErr != nil {
				common.ErrorMsg("InsertProductData", lErr.Error())
				return lErr
			}
		}

		// Order Details Check if record exist if exist not only insert the record into orders table
		lOrderID, lErr := Check_Record_If_Exists("orders", "order_id", lRow[0])
		if lErr != nil {
			common.ErrorMsg("Check_If_RecordExists", lErr.Error())
			return lErr
		}
		if lOrderID != 0 {
			continue
		}
		newOrder := model.OrdersStruct{
			Order_id:       lRow[0],
			Customer_id:    lCustomerID,
			Region:         lRow[5],
			Date_of_sale:   lRow[6],
			Payment_method: lRow[11],
			Shipping_cost:  lRow[10],
		}
		lOrderID, lErr = Insert_Order_Data(newOrder)
		if lErr != nil {
			common.ErrorMsg("Insert_Order_Data", lErr.Error())
			return lErr
		}

		//Here we are inserting the order items list into the order_items table
		lFilteredRecords, lErr := Construct_OrderItem_Records(lRow, lProductID, lOrderID)
		if lErr != nil {
			common.ErrorMsg("Construct_OrderItem_Records", lErr.Error())
			return lErr
		} else {
			lErr = Insert_Order_ItemList(lFilteredRecords)
			if lErr != nil {
				common.ErrorMsg("Insert_Order_ItemList", lErr.Error())
				return lErr
			}
		}
	}
	log.Println("SplitAndInsert (-)")
	return nil
}

func Construct_OrderItem_Records(row []string, orderID, productID int) (model.OrderItemsStruct, error) {
	log.Println("Construct_OrderItem_Records (+)")
	var lOrderedItemlist model.OrderItemsStruct
	var lErr error

	lOrderedItemlist.Order_id = orderID
	lOrderedItemlist.Product_id = productID
	lOrderedItemlist.Quantity_sold, lErr = strconv.Atoi(row[7])
	if lErr != nil {
		return lOrderedItemlist, lErr
	}
	lOrderedItemlist.Unit_price, lErr = strconv.ParseFloat(row[8], 64)
	if lErr != nil {
		return lOrderedItemlist, lErr
	}
	lOrderedItemlist.Discount, lErr = strconv.ParseFloat(row[9], 64)
	if lErr != nil {
		return lOrderedItemlist, lErr
	}
	log.Println("Construct_OrderItem_Records (-)")
	return lOrderedItemlist, nil
}
