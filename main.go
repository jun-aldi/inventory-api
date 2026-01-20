package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Product struct {
	ID    int    `json:"id"`
	Nama  string `json:"nama"`
	Harga int    `json:"harga"`
	Stok  int    `json:"stok"`
}

var produk = []Product{
	{ID: 1, Nama: "Indomie", Harga: 150000, Stok: 10},
	{ID: 2, Nama: "Baju", Harga: 100000, Stok: 20},
	{ID: 3, Nama: "Celana", Harga: 120000, Stok: 15},
}

func getProductByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Produk Id", http.StatusBadRequest)
		return
	}

	for _, p := range produk {
		if p.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(p)
			return
		}
	}

	http.Error(w, "Produk belum ada", http.StatusNotFound)
}

func updateProduk(w http.ResponseWriter, r *http.Request) {
	// GET ID dari request
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")

	// GANTI INT
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Produk Id", http.StatusBadRequest)
		return
	}

	// GET DATA DARI REQUEST

	var updateProduk Product
	err = json.NewDecoder(r.Body).Decode(&updateProduk)
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	// LOOP PRODUK, CARI ID GANTI SESUAI
	for i := range produk {
		updateProduk.ID = id
		produk[i] = updateProduk
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updateProduk)
		return
	}
	http.Error(w, "Produk belum ada", http.StatusNotFound)
}

func deleteProduk(w http.ResponseWriter, r *http.Request) {

	// get id
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")

	// GANTI INT
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Produk Id", http.StatusBadRequest)
		return
	}

	// LOOP PRODUK CARI ID, DAPAT INDEX YANG MAU DIHAPUS
	for i, p := range produk {
		if p.ID == id {

			// BIKIN SLICE BARU DENGAN DATA SEBELUM DAN SESUDAH INDEX
			produk = append(produk[:i], produk[i+1:]...)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "sukses delete",
			})

			return

		}
	}

	http.Error(w, "Produk belum ada", http.StatusNotFound)

}

func main() {

	// GET localhost:8080/api/produk/{id}
	// PUT localhost:8080/api/produk/{id}
	// DELETE localhost:8080/api/produk/{id}
	http.HandleFunc("/api/produk/", func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "GET" {
			getProductByID(w, r)
		} else if r.Method == "PUT" {
			updateProduk(w, r)
		} else if r.Method == "DELETE" {
			deleteProduk(w, r)
		}

	})

	// GET localhost:8080/api/produk
	// POST localhost:8080/api/produk
	http.HandleFunc("/api/produk", func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(produk)
			return
		}

		if r.Method == http.MethodPost {
			var produkBaru Product

			err := json.NewDecoder(r.Body).Decode(&produkBaru)
			if err != nil {
				http.Error(w, "Invalid Request", http.StatusBadRequest)
				return
			}

			// masukan data ke dalam variable produk
			produkBaru.ID = len(produk) + 1
			produk = append(produk, produkBaru)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(produkBaru)
			return
		}

		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	})

	// api/health
	http.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		response := map[string]string{
			"status":  "OK",
			"message": "API running",
		}

		json.NewEncoder(w).Encode(response)
	})

	fmt.Println("Server running di http://localhost:8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("gagal running server:", err)
	}
}
