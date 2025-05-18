package fetchrevenue

import (
	"SalesReport/common"
	"SalesReport/global"
	"SalesReport/model"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func GetCategoryWiseRevenue(w http.ResponseWriter, r *http.Request) {

	log.Println("GetCategoryWiseRevenue (+) ")
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Category, Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, credentials")

	if r.Method == http.MethodPost {

		var lInputRec model.GetDataStruct
		var lRevenueResp model.GetRevenueStruct

		lRevenueResp.Status = common.SuccessCode

		//Read the Header value
		lCategory := r.Header.Get("Category")

		lErr := json.NewDecoder(r.Body).Decode(&lInputRec)
		if lErr != nil {
			lRevenueResp.Status = common.ErrorCode
			common.ErrorMsg("GetCategoryWiseRevenue", lErr.Error())
		} else {
			if strings.TrimSpace(lInputRec.FromDate) == "" || strings.TrimSpace(lInputRec.ToDate) == "" {
				lRevenueResp.Status = common.ErrorCode
			} else {
				switch lCategory {
				case "TotalRevenue":
					lRevenueResp.Total_revenue, lErr = GetOverallRevenue(lInputRec)
				case "Product":
					lRevenueResp.TotProdRevenue, lErr = GetProductsRevenue(lInputRec)
				case "Category":
					lRevenueResp.TotalcatRevenue, lErr = GetCategoryRevenue(lInputRec)
				case "Region":
					lRevenueResp.TotalRevenue_byreg, lErr = GetRegionRevenue(lInputRec)
				default:
					http.Error(w, "Invalid Indicator", http.StatusBadRequest)
				}

				if lErr != nil {
					lRevenueResp.Status = common.ErrorCode
					common.ErrorMsg("GetCategoryWiseRevenue", lErr.Error())
				}

			}
			lData, lErr := json.Marshal(lRevenueResp)
			if lErr != nil {
				common.ErrorMsg("GetCategoryWiseRevenue", lErr.Error())
				fmt.Fprint(w, lErr.Error())
			} else {
				fmt.Fprint(w, string(lData))
			}
		}
	}

	log.Println("GetCategoryWiseRevenue (-) ")
}

func GetOverallRevenue(pInputRec model.GetDataStruct) (string, error) {
	log.Println("GetOverallRevenue (+)")
	//Get COnversion date
	var lOverallRevenue string
	lFromDate := common.GetDate(pInputRec.FromDate)
	lToDate := common.GetDate(pInputRec.ToDate)

	lCoreString := `SELECT COALESCE(product_revenue + shipping_revenue, '0') AS total_revenue
				FROM (
					SELECT
						SUM(oi.quantity_sold * oi.unit_price * (1 - oi.discount)) AS product_revenue,
						SUM(o.shipping_cost) AS shipping_revenue
					FROM orders o
					JOIN order_items oi ON o.id = oi.order_id
					WHERE o.date_of_sale BETWEEN $1 AND $2
				) AS revenue_data`

	log.Println("lCoreString", lCoreString)
	lErr := global.GConnection.DbPostgres.QueryRow(lCoreString, lFromDate, lToDate).Scan(&lOverallRevenue)
	if lErr != nil {
		common.ErrorMsg("GetOverallRevenue", lErr.Error())
		log.Println("Error in GOR01", lErr.Error())
		return "", lErr
	}
	log.Println("GetOverallRevenue (-)")

	return lOverallRevenue, nil
}

func GetProductsRevenue(pInputRec model.GetDataStruct) ([]model.ProductRevenueStruct, error) {
	log.Println(" GetProductsRevenue (+) ")

	var lProductArr []model.ProductRevenueStruct
	var lProductRec model.ProductRevenueStruct

	lFromDate := common.GetDate(pInputRec.FromDate)
	lToDate := common.GetDate(pInputRec.ToDate)

	lCoreString := `SELECT 
				p.product_name AS product_name,
				COALESCE(SUM(oi.quantity_sold * oi.unit_price * (1 - oi.discount)), 0) AS total_revenue
			FROM products p
			JOIN order_items oi ON p.id = oi.product_id
			JOIN orders o ON oi.order_id = o.id
			WHERE o.date_of_sale BETWEEN $1 AND $2			
			GROUP BY p.product_name`

	lRows, lErr := global.GConnection.DbPostgres.Query(lCoreString, lFromDate, lToDate)
	if lErr != nil {
		common.ErrorMsg("GetProductsRevenue", lErr.Error())
		log.Println("Error in GPR01", lErr.Error())
		return nil, lErr
	} else {
		defer lRows.Close()
		for lRows.Next() {
			lErr := lRows.Scan(&lProductRec.ProductName, &lProductRec.TotalRevenue)
			if lErr != nil {
				common.ErrorMsg("GetProductsRevenue", lErr.Error())
				log.Println("Error in GPR02", lErr.Error())
				return nil, lErr
			} else {
				lProductArr = append(lProductArr, lProductRec)
			}
		}
	}

	log.Println(" GetRevenuebyProd (-) ")
	return lProductArr, nil
}

func GetCategoryRevenue(pInputRec model.GetDataStruct) ([]model.CategoryRevenueStruct, error) {
	log.Println("GetCategoryRevenue(+)")

	var lCategoryRec model.CategoryRevenueStruct
	var lCategoryArr []model.CategoryRevenueStruct

	lFromDate := common.GetDate(pInputRec.FromDate)
	lToDate := common.GetDate(pInputRec.ToDate)

	lCoreString := `
					SELECT 
						p.category AS category_name,
						COALESCE(SUM(oi.quantity_sold * oi.unit_price * (1 - oi.discount)), 0) AS total_revenue
					FROM products p
					JOIN order_items oi ON p.id = oi.product_id
					JOIN orders o ON oi.order_id = o.id
					WHERE o.date_of_sale BETWEEN $1 AND $2
					GROUP BY p.category`

	lRows, lErr := global.GConnection.DbPostgres.Query(lCoreString, lFromDate, lToDate)
	if lErr != nil {
		common.ErrorMsg("GetCategoryRevenue", lErr.Error())
		log.Println("Error in GCR01", lErr.Error())
		return nil, lErr
	} else {
		defer lRows.Close()
		for lRows.Next() {
			lErr := lRows.Scan(&lCategoryRec.Category, &lCategoryRec.TotalRevenue)
			if lErr != nil {
				common.ErrorMsg("GetCategoryRevenue", lErr.Error())
				log.Println("Error in GCR02", lErr.Error())
				return nil, lErr
			} else {
				lCategoryArr = append(lCategoryArr, lCategoryRec)
			}
		}
	}
	log.Println("GetCategoryRevenue(-)")
	return lCategoryArr, nil
}

func GetRegionRevenue(pInputRec model.GetDataStruct) ([]model.RegionRevenueStruct, error) {
	log.Println(" GetRegionRevenue (+) ")

	var lRegionArr []model.RegionRevenueStruct
	var lRegionRec model.RegionRevenueStruct
	lFromDate := common.GetDate(pInputRec.FromDate)
	lToDate := common.GetDate(pInputRec.ToDate)

	lCoreString := `SELECT 
						COALESCE(o.region, '-') AS region,
						COALESCE(SUM(oi.quantity_sold * oi.unit_price * (1 - oi.discount)), 0) AS total_revenue
					FROM orders o
					JOIN order_items oi ON o.id = oi.order_id
					WHERE o.date_of_sale BETWEEN $1 AND $2
					GROUP BY o.region`

	lRows, lErr := global.GConnection.DbPostgres.Query(lCoreString, lFromDate, lToDate)
	if lErr != nil {
		common.ErrorMsg("GetRegionRevenue", lErr.Error())
		log.Println("Error in GRR01", lErr.Error())
		return nil, lErr
	} else {
		defer lRows.Close()
		for lRows.Next() {
			lErr := lRows.Scan(&lRegionRec.Region, &lRegionRec.TotalRevenue)
			if lErr != nil {
				common.ErrorMsg("GetRegionRevenue", lErr.Error())
				log.Println("Error in GRR02", lErr.Error())
				return nil, lErr
			} else {
				lRegionArr = append(lRegionArr, lRegionRec)

			}
		}
	}

	log.Println("GetRegionRevenue (-) ")
	return lRegionArr, nil
}
