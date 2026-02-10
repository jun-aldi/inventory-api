package main

import (
	"encoding/json"
	"fmt"
	"kasir-api/database"
	"kasir-api/handlers"
	"kasir-api/middleware"
	"kasir-api/repositories"
	"kasir-api/services"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
	APIKey string `mapstructure:"API_KEY"`
}

func main() {
	// Konfigurasi Environment
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
		APIKey: viper.GetString("API_KEY"),
	}

	//Init Database
	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Setup Middleware & Dependency Injection
	apiKeyMiddleware := middleware.APIKey(config.APIKey)

	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	transactionRepo := repositories.NewTransactionRepository(db)
	transactionService := services.NewTransactionService(transactionRepo)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	// Setup Routes

	// -- Product --
	http.HandleFunc("/api/product", productHandler.HandleProducts)
	http.HandleFunc("/api/product/", middleware.Logger(apiKeyMiddleware(productHandler.HandleProductByID)))

	// -- Category --
	http.HandleFunc("/api/category", categoryHandler.HandleCategories)
	http.HandleFunc("/api/category/", middleware.Logger(apiKeyMiddleware(categoryHandler.HandleCategoryByID)))

	// -- Checkout --
	http.HandleFunc("/api/checkout", middleware.Logger(apiKeyMiddleware(transactionHandler.HandleCheckout)))

	// -- Report --
	http.HandleFunc("/api/report/hari-ini", transactionHandler.HandleReport)
	http.HandleFunc("/api/report", transactionHandler.HandleReport)

	// -- Health Check --
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "API running",
		})
	})

	addr := "0.0.0.0:" + config.Port
	fmt.Println("Server running di", addr)

	globalHandler := middleware.CORS(http.DefaultServeMux)

	err = http.ListenAndServe(addr, globalHandler)
	if err != nil {
		fmt.Println("gagal running server", err)
	}
}
