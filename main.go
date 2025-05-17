package main

import (
	"golang-sqlserver/internal/controllers"
	"golang-sqlserver/internal/database"
	"golang-sqlserver/internal/repository"
	"golang-sqlserver/internal/services"
	"log"
	"net/http"

	_ "github.com/denisenkom/go-mssqldb"
)

func main() {
	// Database configuration
	config := database.DBConfig{
		Server:   "localhost",
		Port:     1433,
		User:     "SA",
		Password: "yourStrong(!)Password",
		Database: "golang",
	}

	db, err := database.NewConnection(config)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer db.Close()

	productRepo := repository.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productController := controllers.NewProductController(productService)

	http.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			productController.Create(w, r)
		case http.MethodGet:
			if r.URL.Query().Get("id") != "" {
				productController.GetByID(w, r)
			} else {
				productController.Search(w, r)
			}
		case http.MethodPut:
			productController.Update(w, r)
		case http.MethodDelete:
			productController.Delete(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/products/search", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		productController.Search(w, r)
	})

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
