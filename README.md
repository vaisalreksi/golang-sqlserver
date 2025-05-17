# golang-sqlserver

Golang SQL Server

# Introduction

Example Golang SQL Server

# Database

```
create database golang
```

# Product
| Field            | Type   | 
| :--------------- |:------:| 
| Id               | int    | 
| Name             | string | 
| Product_Category | string | 
| Description      | string |                               


# Price 
| Field            | Type   | 
| :--------------- |:------:| 
| Id               | int    | 
| Product_Id       | int    | 
| Unit             | string | 


# PriceDetail            
| Field         | Type   |
|:------------- |:------:|
| Id            | int    |
| Price_Id      | int    |
| Tier          | string |
| Price         | float  |

### Running
```

go run *.go
Your Appication is running on port 8080

```
---

### Curl
```

Create Product

curl --location 'http://yourlocal:8080/products' \
--header 'Content-Type: application/json' \
--data '{
"name": "Product ica icaica",
"product_category": "Rokok",
"description": "description"
}'

```

```

Update Product

curl --location --request PUT 'http://yourlocal:8080/products' \
--header 'Content-Type: application/json' \
--data '{
"name": "Product ica icaica",
"product_category": "Rokok",
"description": "description",
"id": 1
}'

```

```

Get Products

curl --location --request GET 'http://yourlocal:8080/products'

```

```

Get Product

curl --location --request GET 'http://yourlocal:8080/products?id=1'

```

```

Delete Product

curl --location --request DELETE 'http://yourlocal:8080/products?id=1'

```

```

Search Product

curl --location 'http://yourlocal:8080/products/search' \
--header 'Content-Type: application/json' \
--data '{
"keyword": "product",
"product_category": "",
"tier": "Basic"
}'

```

