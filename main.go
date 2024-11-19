package main

import (
    "log"
    "net/http"
    "app/controllers"
    _ "app/docs" // Import Swagger docs
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    httpSwagger "github.com/swaggo/http-swagger"
)

// Routes function defines the API routes for the application
func Routes() http.Handler {
    r := chi.NewRouter()
    r.Use(middleware.Logger)

    // Define routes for handling products
    r.Post("/products", controllers.CreateProduct) // create product
    r.Get("/products", controllers.GetAllProducts) // List Products
    r.Get("/products/{id}", controllers.GetProductByID) // Get Product by ID
    r.Put("/products/{id}", controllers.UpdateProduct) // Update Product by ID
    r.Delete("/delete/{id}", controllers.DeleteProduct) // Deletes a Product
    r.Post("/products/many", controllers.BulkInsertProduct) //Create Multiple Products
    r.Delete("/delete/products", controllers.DeleteProducts) // Delete Multiple Products
    r.Get("/products/category/{category}", controllers.CategoryOfProducts) // Gets Products by Category

    // Swagger documentation route
    r.Get("/swagger/*", httpSwagger.WrapHandler)

    return r
}

func main() {
    // Initialize MongoDB connection
    controllers.InitMongoDB()

    // Setup routes
    r := Routes()

    // Start the server
    log.Println("Starting server on :8080...")
    if err := http.ListenAndServe(":8080", r); err != nil {
        log.Fatal(err)
    }
}
