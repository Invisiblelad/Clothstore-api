# Product Management API

This API allows you to manage products with various endpoints, including operations for creating, updating, retrieving, and deleting products. Below are the available API routes along with their respective `curl` commands.

## MongoDB

This API uses MongoDB as the backend to store product data. And the MongoDb is run as a container in your local machine.

```bash
docker run --name mongodb -d -p 27017:27017 mongo
```

### Run the code using go run command

```bash 

go run main.go

```

## Endpoints

### POST /products - Create a New Product

Create a new product by sending a `POST` request with a JSON body.

#### Request:

```bash
curl -X POST http://localhost:8080/products \
-H "Content-Type: application/json" \
-d '{
    "name": "Shirt",
    "category": "Clothing",
    "price": 29.99,
    "description": "A stylish shirt for all occasions"
}'

```

### GET /products - Get All Products

Retrieve all products by sending a GET request.

```bash
curl -X GET http://localhost:8080/products
```

### GET /products/{id} - Get Product by ID

Retrieve a single product by its unique ID by sending a GET request.

```bash
curl -X GET http://localhost:8080/products/<product_id>

```

### PUT /products/{id} - Update Product

Update an existing product by sending a PUT request with the updated product data.

```bash
curl -X PUT http://localhost:8080/products/<product_id> \
-H "Content-Type: application/json" \
-d '{
    "name": "T-Shirt",
    "category": "Clothing",
    "price": 24.99,
    "description": "An updated description of the product"
}'

```

### DELETE /delete/{id} - Delete Product by ID

Delete a product by its unique ID by sending a DELETE request.

```bash
curl -X DELETE http://localhost:8080/delete/<product_id>

```

### POST /products/many - Bulk Insert Products

Insert multiple products at once by sending a POST request with a JSON array of products.

```bash
curl -X POST http://localhost:8080/products/many \
-H "Content-Type: application/json" \
-d '[ 
    {
        "name": "T-Shirt",
        "category": "Clothing",
        "price": 19.99,
        "description": "A casual t-shirt"
    },
    {
        "name": "Jeans",
        "category": "Clothing",
        "price": 49.99,
        "description": "Comfortable denim jeans"
    }
]'

```
### DELETE /delete/products - Bulk Delete Products by Category

Delete multiple products based on a category by sending a DELETE request with the category name in the request body.

```bash

curl -X DELETE http://localhost:8080/delete/products \
-H "Content-Type: application/json" \
-d '{
    "category": "Clothing"
}'

```
### Get /products/category/{category} - Fetch product by Category

Fetch the list of products in the specified category.

```bash

curl -X GET http://localhost:8080/products/category/<category>

```

