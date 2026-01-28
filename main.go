package main

import (
	"encoding/json"
	"fmt"
	"kasir-api/database"
	"kasir-api/handlers"
	"kasir-api/repositories"
	"kasir-api/services"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

// categories struct
type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

var categories = []Category{
	{ID: 1, Name: "Raw Materials", Description: "Basic materials used as inputs in the production process"},
	{ID: 2, Name: "Work in Process", Description: "Partially completed goods that are still in production"},
	{ID: 3, Name: "Finished Goods", Description: "Final products that are ready for sale or distribution"},
}

var lastID int

// INIT LAST ID
func init() {
	for _, c := range categories {
		if c.ID > lastID {
			lastID = c.ID
		}
	}
}

// GET DETAIL CATEGORY
func getCategoryByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Category Id", http.StatusBadRequest)
		return
	}

	for _, p := range categories {
		if p.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(p)
			return
		}
	}

	http.Error(w, "Category not found", http.StatusNotFound)
}

// PUT UPDATE CATEGORY
func updateCategory(w http.ResponseWriter, r *http.Request) {
	// GET ID dari request
	idStr := strings.TrimPrefix(r.URL.Path, "/categories/")

	// GANTI INT
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Category Id", http.StatusBadRequest)
		return
	}

	// GET DATA DARI REQUEST

	var updateCategory Category
	err = json.NewDecoder(r.Body).Decode(&updateCategory)
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	// LOOP CATEGORY, CARI ID GANTI SESUAI
	for i := range categories {
		if categories[i].ID == id {
			categories[i].Name = updateCategory.Name
			categories[i].Description = updateCategory.Description

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(categories[i])
			return
		}
	}

	http.Error(w, "Category not found", http.StatusNotFound)
}

// DELETE CATEGORY
func deleteCategory(w http.ResponseWriter, r *http.Request) {

	// get id
	idStr := strings.TrimPrefix(r.URL.Path, "/categories/")

	// GANTI INT
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Category Id"+err.Error(), http.StatusBadRequest)
		return
	}

	// LOOP CATEGORIES CARI ID, DAPAT INDEX YANG MAU DIHAPUS
	for i, p := range categories {
		if p.ID == id {

			// BIKIN SLICE BARU DENGAN DATA SEBELUM DAN SESUDAH INDEX
			categories = append(categories[:i], categories[i+1:]...)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "delete success",
			})

			return

		}
	}

	http.Error(w, "Category not found", http.StatusNotFound)

}

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

func main() {

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	// Setup routes
	http.HandleFunc("/api/produk", productHandler.HandleProducts)
	http.HandleFunc("/api/produk/", productHandler.HandleProductByID)

	http.HandleFunc("/categories/", func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "GET" {
			getCategoryByID(w, r)
		} else if r.Method == "PUT" {
			updateCategory(w, r)
		} else if r.Method == "DELETE" {
			deleteCategory(w, r)
		}

	})

	// GET localhost:8080/category
	// POST localhost:8080/category
	http.HandleFunc("/categories", func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(categories)
			return
		}

		if r.Method == http.MethodPost {

			var newCategory Category
			if err := json.NewDecoder(r.Body).Decode(&newCategory); err != nil {
				http.Error(w, "Invalid Request", http.StatusBadRequest)
				return
			}

			lastID++
			newCategory.ID = lastID

			categories = append(categories, newCategory)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(newCategory)
			return
		}

		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	})

	// /health
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		response := map[string]string{
			"status":  "OK",
			"message": "API running",
		}

		json.NewEncoder(w).Encode(response)
	})

	addr := "0.0.0.0:" + config.Port
	fmt.Println("Server running di", addr)

	err = http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println("gagal running server", err)
	}

}
