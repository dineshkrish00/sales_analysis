# Sales Report Project

This project provides an API for processing and managing large sales datasets. The application allows users to upload CSV files containing historical sales data and perform operations to retrieve sales revenue and product details. The project is built using **Go (Golang)**, **PostgreSQL**, and includes proper error handling.

---
## Git Clone URL
```
git@github.com:dineshkrish00/sales_report.git
```

## Install Dependencies

go mod tidy

## Run

go run main.go

## LOG file 
The application creates logs while running to help track the execution flow and debug issues.
Log Location: log/logfile19042025.12.20.03.610924298.txt

## Technologies Used

- **Go (Golang)**: Backend logic, CSV file processing, API routing
- **PostgreSQL**: Database for storing sales records
- **JSON API**: For communication between frontend and backend

---

## API Endpoints

### 1. **POST /csvfileupload**

**URL**: `http://localhost:29001/csvfileupload`

This endpoint allows you to upload a CSV file containing historical sales data. The file will be processed and its data inserted into the PostgreSQL database.

**Request**:
- **Method**: POST
- **Content-Type**: multipart/form-data
- **Body**: File object attached to the request.

If no file is provided in the request, the server will attempt to fetch the file from a specified file path and process it.

**Response**:
- **200 OK**:
  ```json
    {
    "status": "S",
    "msg": "File Uploaded Successfully"
    }
- **status as E:**:
 ```json
    {
    "status": "E",
    "msg": "Error message"
    }
```
---

### 2. **POST /fetchRevenue**

**URL**: `http://localhost:29001/fetchRevenue`

This endpoint retrieves sales revenue data based on the provided date range and category.


**Request**:
- **Method**: POST
- **Content-Type**: multipart/form-data
- **Body**: JSON object containing the fromDate and toDate parameters.
- **Header**: Category: Value can be one of the following: TotalRevenue, Product, Category, Region.
---

**Response**:
- **200 OK**:
- 
###  Response 1:

  ```json
        {
    "status": "S",
    "total_Revenue": "44.5675",
    "totProdRevenue": null,
    "totalcatRevenue": null,
    "totalregionRevenue": null,
    "topProduct": null,
    "topcategory": null,
    "topRegion": null
}
```

###  Response 2:

  ```json
   {
    "status": "S",
    "total_Revenue": "",
    "totProdRevenue": null,
    "totalcatRevenue": [
        {
            "category": "",
            "total_revenue": 143.976
        },
        {
            "category": "",
            "total_revenue": 4064.55
        },
        {
            "category": "",
            "total_revenue": 183
        }
    ],
    "totalregionRevenue": null,
    "topProduct": null,
    "topcategory": null,
    "topRegion": null
}

```

- **status as E:**:
 ```json
   {
    "status": "E",
    "statusCode": "GOR01",
    "msg": "EOF"
}
```
---

### 3. **POST /fetchRevenue**

**URL**: `http://localhost:29001/fetchProductDetails`

This endpoint fetches details about products, based on the given date range and the type of report requested.


**Request**:
- **Method**: POST
- **Content-Type**: multipart/form-data
- **Body**: JSON object containing the fromDate and toDate parameters.
- **Header**: Type: Value can be one of the following: Overall, Category, Region.
---

**Response**:
- **200 OK**:
  ```json
    {
        "status": "S",
        "total_Revenue": "",
        "totProdRevenue": null,
        "totalcatRevenue": null,
        "totalregionRevenue": null,
        "topProduct": [
            {
                "product_name": "",
                "total_quantity": 1
            },
            {
                "product_name": "",
                "total_quantity": 2
            },
            {
                "product_name": "",
                "total_quantity": 1
            },
            {
                "product_name": "",
                "total_quantity": 1
            }
        ],
        "topcategory": null,
        "topRegion": null
    }

- **status as E:**:
 ```json
    {
    "status": "E",
    "msg": "Error message"
    }
```