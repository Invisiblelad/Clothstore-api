package main

import (
    "net/http"

    "cloth-store/controllers"
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
)

func Routes() http.Handler {
    r := chi.NewRouter()
    r.Use(middleware.Logger)

    r.Post("/products", controllers.CreateProduct)         // Create product
    r.Get("/products", controllers.GetAllProducts)         // Get all products
    r.Get("/products/{id}", controllers.GetProductByID)    // Get product by ID

    return r
}

