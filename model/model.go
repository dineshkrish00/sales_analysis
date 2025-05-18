package model

type ResponseStruct struct {
	Status string `json:"status"`
	Msg    string `json:"msg"`
}

type CustomerDataStruct struct {
	ID               int    `json:"id"`
	Customer_id      string `json:"customer_id"`
	Customer_name    string `json:"name"`
	Customer_email   string `json:"email"`
	Customer_address string `json:"address"`
	Created_By       string `json:"created_by"`
	Created_Date     string `json:"created_date"`
	Updated_By       string `json:"updated_by"`
	Updated_Date     string `json:"updated_date"`
}

type ProductsStruct struct {
	ID           int    `json:"id"`
	Product_id   string `json:"product_id"`
	Product_Name string `json:"product_name"`
	Category     string `json:"category"`
	Created_By   string `json:"created_by"`
	Created_Date string `json:"created_date"`
	Updated_By   string `json:"updated_by"`
	Updated_Date string `json:"updated_date"`
}

type OrdersStruct struct {
	ID             int    `json:"id"`
	Order_id       string `json:"order_id"`
	Customer_id    int    `json:"customer_id"`
	Region         string `json:"region"`
	Date_of_sale   string `json:"date_of_sale"`
	Payment_method string `json:"payment_method"`
	Shipping_cost  string `json:"shipping_cost"`
	Created_By     string `json:"created_by"`
	Created_Date   string `json:"created_date"`
	Updated_By     string `json:"updated_by"`
	Updated_Date   string `json:"updated_date"`
}

type OrderItemsStruct struct {
	Order_id      int     `json:"order_id"`
	Product_id    int     `json:"product_id"`
	Quantity_sold int     `json:"quantity_sold"`
	Unit_price    float64 `json:"unit_price"`
	Discount      float64 `json:"discount"`
	Created_By    string  `json:"created_by"`
	Created_Date  string  `json:"created_date"`
	Updated_By    string  `json:"updated_by"`
	Updated_Date  string  `json:"updated_date"`
}

type GetDataStruct struct {
	FromDate string `json:"fromDate"`
	ToDate   string `json:"toDate"`
}

type TopProductStruct struct {
	ProductName   string `json:"product_name"`
	TotalQuantity int    `json:"total_quantity"`
}

type TopCategoryStruct struct {
	Category      string `json:"category"`
	TotalQuantity int    `json:"total_quantity"`
}

type TopRegionStruct struct {
	Region        string `json:"region"`
	TotalQuantity int    `json:"total_quantity"`
}
type ProductRevenueStruct struct {
	ProductName  string  `json:"product_name"`
	TotalRevenue float64 `json:"total_revenue"`
}

type CategoryRevenueStruct struct {
	Category     string  `json:"category"`
	TotalRevenue float64 `json:"total_revenue"`
}

type RegionRevenueStruct struct {
	Region       string  `json:"region"`
	TotalRevenue float64 `json:"total_revenue"`
}
type GetRevenueStruct struct {
	Status             string                  `json:"status"`
	Total_revenue      string                  `json:"total_Revenue"`
	TotProdRevenue     []ProductRevenueStruct  `json:"totProdRevenue"`
	TotalcatRevenue    []CategoryRevenueStruct `json:"totalcatRevenue"`
	TotalRevenue_byreg []RegionRevenueStruct   `json:"totalregionRevenue"`
	TopProduct         []TopProductStruct      `json:"topProduct"`
	TopCategory        []TopCategoryStruct     `json:"topcategory"`
	TopRegion          []TopRegionStruct       `json:"topRegion"`
}
