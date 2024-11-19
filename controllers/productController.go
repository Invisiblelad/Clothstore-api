package controllers

import (
    "context"
    "encoding/json"
    "net/http"

    "app/models"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/bson/primitive" 
    "github.com/go-chi/chi/v5"  
)

var productCollection *mongo.Collection

func InitMongoDB() {
    clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
    client, err := mongo.Connect(context.Background(), clientOptions)
    if err != nil {
        panic(err)
    }
    err = client.Ping(context.Background(), nil)
    if err != nil {
        panic(err)
    }
    productCollection = client.Database("clothstore").Collection("products")
}


// CreateProduct creates a new product in the database
// @Summary Create a product
// @Description Create a new product with the given details
// @Tags products
// @Accept json
// @Produce json
// @Param product body models.Product true "Product Data"
// @Success 201 {object} models.Product
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /products [post]
func CreateProduct(w http.ResponseWriter, r *http.Request) {
    var product models.Product
    err := json.NewDecoder(r.Body).Decode(&product)
    if err != nil {
        http.Error(w, "Error parsing request", http.StatusBadRequest)
        return
    }
    var existingProduct models.Product
    err = productCollection.FindOne(context.Background(), bson.M{"name": product.Name}).Decode(&existingProduct)
    if err == nil {
        http.Error(w, "Product already exists", http.StatusConflict)
        return
    }

    product.ID = primitive.NewObjectID()
    _, err = productCollection.InsertOne(context.Background(), product)
    if err != nil {
        http.Error(w, "Failed to insert product", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(product)
}

// GetAllProducts retrieves all products from the database
// GetAllProducts retrieves all products from the database
// @Summary Get all products
// @Description Retrieve all products from the database
// @Tags products
// @Accept json
// @Produce json
// @Success 200 {array} models.Product
// @Failure 500 {string} string "Internal Server Error"
// @Router /products [get]
func GetAllProducts(w http.ResponseWriter, r *http.Request) {
    cursor, err := productCollection.Find(context.Background(), bson.D{{}})
    if err != nil {
        http.Error(w, "Failed to fetch products", http.StatusInternalServerError)
        return
    }
    var products []models.Product
    if err := cursor.All(context.Background(), &products); err != nil {
        http.Error(w, "Failed to decode products", http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(products)
}



// GetProductByID retrieves a product by its ID
// @Summary Get a product by ID
// @Description Retrieve a specific product using its unique ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} models.Product
// @Failure 400 {string} string "Invalid Product ID"
// @Failure 404 {string} string "Product Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /products/{id} [get]
func GetProductByID(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    objectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        http.Error(w, "Invalid product ID", http.StatusBadRequest)
        return
    }

    var product models.Product
    err = productCollection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&product)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            http.Error(w, "Product not found", http.StatusNotFound)
        } else {
            http.Error(w, "Failed to retrieve product", http.StatusInternalServerError)
        }
        return
    }

    json.NewEncoder(w).Encode(product)
}

// UpdateProduct updates an existing product in the database
// @Summary Update a product
// @Description Update the details of an existing product
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param product body models.Product true "Updated Product Data"
// @Success 200 {object} map[string]string "Product updated successfully"
// @Failure 400 {string} string "Invalid Product ID or Request Body"
// @Failure 404 {string} string "Product Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /products/{id} [put]
func UpdateProduct(w http.ResponseWriter, r *http.Request){
    id := chi.URLParam(r, "id")
    objectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        http.Error(w, "Invalid product ID", http.StatusBadRequest)
        return
    }
    var updatedproduct models.Product
    err = json.NewDecoder(r.Body).Decode(&updatedproduct)
    if err != nil{
        http.Error(w, "Error parsing request body", http.StatusBadRequest)
        return
    }
    update := bson.M{
        "$set": bson.M{
            "name":        updatedproduct.Name,
            "category": updatedproduct.Category,
            "price":       updatedproduct.Price,
            "description":    updatedproduct.Description,
        },
    }

    result, err := productCollection.UpdateOne(context.Background(), bson.M{"_id": objectID}, update)
    if err != nil {
        http.Error(w, "Failed to update product", http.StatusInternalServerError)
        return
    }

    if result.MatchedCount == 0 {
        http.Error(w, "Product not found", http.StatusNotFound)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Product updated successfully"})
}

// DeleteProduct deletes a product by its ID
// @Summary Delete a product
// @Description Delete a specific product by its unique ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} map[string]string "Product deleted successfully"
// @Failure 400 {string} string "Invalid Product ID"
// @Failure 404 {string} string "Product Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /products/{id} [delete]
func DeleteProduct(w http.ResponseWriter, r *http.Request){
    id := chi.URLParam(r, "id")
    objectID, err := primitive.ObjectIDFromHex(id) 
    if err != nil {
        http.Error(w, "Invalid product ID", http.StatusBadRequest)
        return
    }

    result, err := productCollection.DeleteOne(context.Background(), bson.M{"_id": objectID})

    if err != nil{
        http.Error(w,"Failed to delete product ", http.StatusInternalServerError)
        return
    }

    if result.DeletedCount==0{
        http.Error(w,"Product not found", http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(map[string]string{"message": "Product deleted Successfully"})

}


// BulkInsertProduct inserts multiple products into the database
// @Summary Bulk insert products
// @Description Insert multiple products into the database at once
// @Tags products
// @Accept json
// @Produce json
// @Param products body []models.Product true "List of Products"
// @Success 201 {object} map[string]interface{} "Products uploaded successfully"
// @Failure 400 {string} string "Invalid Input"
// @Failure 500 {string} string "Internal Server Error"
// @Router /products/bulk [post]
func BulkInsertProduct(w http.ResponseWriter, r *http.Request){
    var products []models.Product
    err := json.NewDecoder(r.Body).Decode(&products)

    if err != nil{
        http.Error(w, "Invalid Input" , http.StatusBadRequest)
        return
    }
    var documents []interface{}
    for _,product := range products {
        documents = append(documents, product)
    }

    result , err := productCollection.InsertMany(context.Background(), documents)

    if err != nil {
        http.Error(w, "Failed to insert products" , http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(map[string]interface{}{
        "message": "Products uploaded sucessfully",
        "insert_ids": result.InsertedIDs,
    })
}

// DeleteProducts deletes multiple products based on a filter
// @Summary Delete multiple products
// @Description Delete products matching the specified filter criteria
// @Tags products
// @Accept json
// @Produce json
// @Param filter body map[string]interface{} true "Filter Criteria"
// @Success 200 {object} map[string]interface{} "Products deleted successfully"
// @Failure 400 {string} string "Invalid Filter Criteria"
// @Failure 500 {string} string "Internal Server Error"
// @Router /products [delete]
func DeleteProducts(w http.ResponseWriter, r *http.Request){
    var filter map[string]interface{}
    if err := json.NewDecoder(r.Body).Decode(&filter); err != nil{
        http.Error(w, "invalid filter criteria", http.StatusBadRequest)
        return
    }
    
    result, err := productCollection.DeleteMany(context.Background(), filter)

    if err != nil {
        http.Error(w, "Failed to delete products", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "message": "Products deleted successfully",
        "deleted_count": result.DeletedCount,
})
}

// CategoryOfProducts godoc
// @Summary Get products by category
// @Description Retrieve a list of products for a given category
// @Tags products
// @Accept json
// @Produce json
// @Param category path string true "Product Category" example(electronics)
// @Success 200 {array} models.Product "List of products"
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Category not found"
// @Failure 500 {string} string "Internal server error"
// @Router /products/category/{category} [get]
func CategoryOfProducts(w http.ResponseWriter, r *http.Request) {
    category := chi.URLParam(r, "category")
    
    var products []models.Product
    cursor, err := productCollection.Find(context.Background(), bson.M{"category": category})
    if err != nil {
        if err == mongo.ErrNoDocuments {
            http.Error(w, "Category not found", http.StatusNotFound)
        } else {
            http.Error(w, "Failed to retrieve products by category", http.StatusInternalServerError)
        }
        return
    }

    if err := cursor.All(context.Background(), &products); err != nil {
        http.Error(w, "Failed to decode products", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(products)
}

// ReplaceProduct replaces an existing product in the database
// @Summary Replace a product
// @Description Replace the details of an existing product with the provided new details
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param product body models.Product true "Product Data"
// @Success 200 {object} models.Product "Product replaced successfully"
// @Failure 400 {string} string "Invalid Product ID or Request Body"
// @Failure 404 {string} string "Product Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /products/{id}/replace [put]
func ReplaceProduct(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    objectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        http.Error(w, "Invalid product ID", http.StatusBadRequest)
        return
    }

    var product models.Product
    err = json.NewDecoder(r.Body).Decode(&product)
    if err != nil {
        http.Error(w, "Error parsing request body", http.StatusBadRequest)
        return
    }

    result, err := productCollection.ReplaceOne(
        context.Background(),
        bson.M{"_id": objectID}, 
        product,                 
    )
    if err != nil {
        http.Error(w, "Failed to replace product", http.StatusInternalServerError)
        return
    }

    if result.MatchedCount == 0 {
        http.Error(w, "Product not found", http.StatusNotFound)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(product)
}







