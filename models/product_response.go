package models

type ProductResponse struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Price        float64 `json:"price"`
	Stock        int     `json:"stock"`
	CategoryID   int     `json:"category_id"`
	CategoryName string  `json:"category_name"`
}
