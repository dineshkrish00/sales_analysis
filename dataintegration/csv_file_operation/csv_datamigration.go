package fileoperation

import (
	"SalesReport/common"
	"SalesReport/global"
	"SalesReport/model"
	"database/sql"
	"fmt"
	"log"
	"strconv"
)

func Check_Record_If_Exists(pTableName, pColumn, pInputval string) (int, error) {
	log.Println("Check_Record_If_Exists (+)")

	var lRecordId int

	lSqlquery := fmt.Sprintf("SELECT id FROM %s WHERE %s = $1", pTableName, pColumn)
	// log.Println("lSqlquery: ", lSqlquery)

	lErr := global.GConnection.DbPostgres.QueryRow(lSqlquery, pInputval).Scan(&lRecordId)
	if lErr == sql.ErrNoRows {
		common.ErrorMsg("Check_Record_If_Exists", lErr.Error())
		return lRecordId, nil
	}

	log.Println("Check_Record_If_Exists (-)")
	return lRecordId, nil
}

func Insert_Customer_Data(pInputCustomer model.CustomerDataStruct) (int, error) {
	log.Println("Insert_Customer_Data (+)")

	var lRecordId int

	lCorestring := `insert into customers (customer_id, customer_name, customer_email, customer_address, created_by,created_date,updated_by,updated_date)
			  VALUES (` + pInputCustomer.Customer_id + `,'` + pInputCustomer.Customer_name + `','` + pInputCustomer.Customer_email + `', '` + pInputCustomer.Customer_address + `','AUTOBOT',NOW(),'AUTOBOT',NOW()) RETURNING id`

	lErr := global.GConnection.DbPostgres.QueryRow(lCorestring).Scan(&lRecordId)
	if lErr != nil {
		common.ErrorMsg("Insert_Customer_Data", lErr.Error())
		return lRecordId, lErr
	}

	log.Println("Insert_Customer_Data (-)")
	return lRecordId, nil
}

func Insert_Product_Data(pInputProduct model.ProductsStruct) (int, error) {
	log.Println("Insert_Product_Data (+)")

	var lRecordId int

	lCorestring := `insert into products (product_id, product_name, category, created_by,created_date,updated_by,updated_date)
			  VALUES (` + pInputProduct.Product_id + `,'` + pInputProduct.Product_Name + `','` + pInputProduct.Category + `','AUTOBOT',NOW(),'AUTOBOT',NOW()) RETURNING id`

	lErr := global.GConnection.DbPostgres.QueryRow(lCorestring).Scan(&lRecordId)
	if lErr != nil {
		common.ErrorMsg("Insert_Product_Data", lErr.Error())
		return lRecordId, lErr
	}

	log.Println("Insert_Product_Data (-)")
	return lRecordId, nil
}

func Insert_Order_Data(pInputOrder model.OrdersStruct) (int, error) {
	log.Println("Insert_Order_Data (+)")

	var lRecordId int

	lCorestring := `insert into orders (order_id, customer_id, region, date_of_sale, payment_method, shipping_cost, created_by,created_date,updated_by,updated_date)
			  VALUES ('` + pInputOrder.Order_id + `', '` + strconv.Itoa(pInputOrder.Customer_id) + `','` + pInputOrder.Region + `','` + pInputOrder.Date_of_sale + `', '` + pInputOrder.Payment_method + `',` + pInputOrder.Shipping_cost + `,'AUTOBOT',NOW(),'AUTOBOT',NOW()) RETURNING id`

	lErr := global.GConnection.DbPostgres.QueryRow(lCorestring).Scan(&lRecordId)
	if lErr != nil {
		common.ErrorMsg("Insert_Order_Data", lErr.Error())
		return lRecordId, lErr
	}

	log.Println("Insert_Order_Data (-)")
	return lRecordId, nil
}

func Insert_Order_ItemList(pInputItem model.OrderItemsStruct) error {
	log.Println("Insert_Order_ItemList (+)")

	lCorestring := `insert into order_items (order_id, product_id, quantity_sold, unit_price, discount, created_by,created_date,updated_by,updated_date)
			  VALUES ($1, $2, $3, $4, $5,'AUTOBOT',NOW(),'AUTOBOT',NOW())`

	_, lErr := global.GConnection.DbPostgres.Exec(lCorestring, pInputItem.Order_id, pInputItem.Product_id, pInputItem.Quantity_sold,
		pInputItem.Unit_price, pInputItem.Discount)
	if lErr != nil {
		common.ErrorMsg("Insert_Order_ItemList", lErr.Error())
		return lErr
	}

	log.Println("Insert_Order_ItemList (-)")
	return nil
}
