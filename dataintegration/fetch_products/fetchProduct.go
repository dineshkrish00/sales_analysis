package fetchproducts

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

func GetCategoryWiseProduct(w http.ResponseWriter, r *http.Request) {
	log.Println(" GetCategoryWiseProduct (+) ")

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Type, Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, credentials")

	if r.Method == http.MethodPost {

		var lInputRec model.GetDataStruct
		var lRespRec model.GetRevenueStruct
		var lErr error

		lType := r.Header.Get("Type")

		lRespRec.Status = "S" // Success

		lErr = json.NewDecoder(r.Body).Decode(&lInputRec)
		if lErr != nil {
			lRespRec.Status = common.ErrorCode
			common.ErrorMsg("GetCategoryWiseProduct", lErr.Error())
		} else {

			if strings.TrimSpace(lInputRec.FromDate) == "" || strings.TrimSpace(lInputRec.ToDate) == "" {
				lRespRec.Status = common.ErrorCode
				log.Println("Empty FIeld is not allowed")
			} else {
				switch lType {

				case "Overall":
					lRespRec.TopProduct, lErr = FetchOverallProducts(lInputRec)
				case "Category":
					lRespRec.TopCategory, lErr = FetchTopCategories(lInputRec)
				case "Region":
					lRespRec.TopRegion, lErr = FetchTopRegions(lInputRec)
				default:
					http.Error(w, "Invalid Indicator", http.StatusBadRequest)

				}

				if lErr != nil {
					lRespRec.Status = common.ErrorCode
					common.ErrorMsg("GetCategoryWiseProduct", lErr.Error())
				}

				lData, lErr := json.Marshal(lRespRec)
				if lErr != nil {
					common.ErrorMsg("GetCategoryWiseProduct", lErr.Error())
					fmt.Fprint(w, lErr.Error())
				} else {
					fmt.Fprint(w, string(lData))
				}
			}

		}

	}
	log.Println(" GetCategoryWiseProduct (-) ")
}

func FetchOverallProducts(input model.GetDataStruct) ([]model.TopProductStruct, error) {
	log.Println("FetchOverallProducts (+)")
	var lProductArr []model.TopProductStruct
	var lProductRec model.TopProductStruct

	lCoreString := `
					SELECT 
					COALESCE(p.product_name, '-') AS product_name, 
					COALESCE(SUM(oi.quantity_sold), 0) AS total_quantity
						FROM order_items oi
						JOIN products p ON oi.product_id = p.id
						JOIN orders o ON oi.order_id = o.id
						WHERE o.date_of_sale BETWEEN $1 AND $2
						GROUP BY p.product_name
						ORDER BY total_quantity DESC`

	lRows, lErr := global.GConnection.DbPostgres.Query(lCoreString, input.FromDate, input.ToDate)
	if lErr != nil {
		common.ErrorMsg("FetchOverallProducts", lErr.Error())
		log.Println("Error in FOP01", lErr.Error())
		return nil, lErr
	} else {
		defer lRows.Close()
		for lRows.Next() {
			lErr := lRows.Scan(&lProductRec.ProductName, &lProductRec.TotalQuantity)
			if lErr != nil {
				common.ErrorMsg("FetchOverallProducts", lErr.Error())
				log.Println("Error in FOP02", lErr.Error())
				return nil, lErr
			} else {
				lProductArr = append(lProductArr, lProductRec)
			}
		}
	}
	return lProductArr, nil
}

func FetchTopCategories(input model.GetDataStruct) ([]model.TopCategoryStruct, error) {
	log.Println("FetchTopCategories (+)")

	var lCategoryRec model.TopCategoryStruct
	var lCategoryArr []model.TopCategoryStruct
	lCoreString := `
		SELECT 
			COALESCE(p.category, '-') AS category, 
			COALESCE(SUM(oi.quantity_sold), 0) AS total_quantity
		FROM order_items oi
		JOIN products p ON oi.product_id = p.id
		JOIN orders o ON oi.order_id = o.id
		WHERE o.date_of_sale BETWEEN $1 AND $2
		GROUP BY p.category
		ORDER BY total_quantity DESC`

	lRows, lErr := global.GConnection.DbPostgres.Query(lCoreString, input.FromDate, input.ToDate)
	if lErr != nil {
		common.ErrorMsg("FetchTopCategories", lErr.Error())
		log.Println("Error in FTC01", lErr.Error())
		return lCategoryArr, lErr
	} else {
		defer lRows.Close()
		for lRows.Next() {
			lErr := lRows.Scan(&lCategoryRec.Category, &lCategoryRec.TotalQuantity)
			if lErr != nil {
				common.ErrorMsg("FetchTopCategories", lErr.Error())
				log.Println("Error in FTC02", lErr.Error())
				return lCategoryArr, lErr
			} else {
				lCategoryArr = append(lCategoryArr, lCategoryRec)
			}
		}
	}
	log.Println("FetchTopCategories (-)")
	return lCategoryArr, nil
}

func FetchTopRegions(input model.GetDataStruct) ([]model.TopRegionStruct, error) {
	log.Println("FetchTopRegions (+)")

	var lRegionsArr []model.TopRegionStruct
	var lRegionsRec model.TopRegionStruct
	lCoreString := `
		SELECT 
			COALESCE(o.region, '-') AS region, 
			COALESCE(SUM(oi.quantity_sold), 0) AS total_quantity
		FROM order_items oi
		JOIN orders o ON oi.order_id = o.id
		WHERE o.date_of_sale BETWEEN $1 AND $2
		GROUP BY o.region
		ORDER BY total_quantity DESC
	`
	lRows, lErr := global.GConnection.DbPostgres.Query(lCoreString, input.FromDate, input.ToDate)
	if lErr != nil {
		common.ErrorMsg("FetchOverallProducts", lErr.Error())
		log.Println("Error in FOP01", lErr.Error())
		return lRegionsArr, lErr
	} else {
		defer lRows.Close()
		for lRows.Next() {
			lErr := lRows.Scan(&lRegionsRec.Region, &lRegionsRec.TotalQuantity)
			if lErr != nil {
				common.ErrorMsg("FetchOverallProducts", lErr.Error())
				log.Println("Error in FOP02", lErr.Error())
				return lRegionsArr, lErr
			} else {
				lRegionsArr = append(lRegionsArr, lRegionsRec)
			}
		}
	}
	log.Println("FetchTopRegions (-)")
	return lRegionsArr, nil
}
