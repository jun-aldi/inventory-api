package models

type BestSellingProduct struct {
	Name      string `json:"nama"`
	TotalSold int    `json:"qty_terjual"`
}

type SalesReport struct {
	TotalRevenue     int                `json:"total_revenue"`
	TotalTransaction int                `json:"total_transaksi"`
	TopProduct       BestSellingProduct `json:"produk_terlaris"`
}
