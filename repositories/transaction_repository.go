package repositories

import (
	"database/sql"
	"fmt"
	"kasir-api/models"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (repo *TransactionRepository) CreateTransaction(items []models.CheckoutItem) (*models.Transaction, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	totalAmount := 0
	details := make([]models.TransactionDetail, 0)

	for _, item := range items {
		var productPrice, stock int
		var productName string

		err := tx.QueryRow("SELECT name, price, stock FROM products WHERE id = $1", item.ProductID).Scan(&productName, &productPrice, &stock)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product id %d not found", item.ProductID)
		}
		if err != nil {
			return nil, err
		}

		if stock < item.Quantity {
			return nil, fmt.Errorf("insufficient stock for product %s", productName)
		}

		subtotal := productPrice * item.Quantity
		totalAmount += subtotal

		_, err = tx.Exec("UPDATE products SET stock = stock - $1 WHERE id = $2", item.Quantity, item.ProductID)
		if err != nil {
			return nil, err
		}

		details = append(details, models.TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	var transactionID int
	err = tx.QueryRow("INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id", totalAmount).Scan(&transactionID)
	if err != nil {
		return nil, err
	}

	if len(details) > 0 {
		query := "INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES "
		var args []interface{}

		for i := range details {
			details[i].TransactionID = transactionID
			base := i * 4
			query += fmt.Sprintf("($%d, $%d, $%d, $%d),", base+1, base+2, base+3, base+4)
			args = append(args, transactionID, details[i].ProductID, details[i].Quantity, details[i].Subtotal)
		}

		query = query[:len(query)-1]

		_, err = tx.Exec(query, args...)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		Details:     details,
	}, nil
}

// SALES REPORT

// Tambahkan method ini di struct TransactionRepository

func (repo *TransactionRepository) GetSalesReport(startDate, endDate string) (*models.SalesReport, error) {
	report := &models.SalesReport{}

	queryStat := `
		SELECT 
			COALESCE(SUM(total_amount), 0), 
			COUNT(id) 
		FROM transactions 
		WHERE created_at >= $1 AND created_at <= $2`

	err := repo.db.QueryRow(queryStat, startDate, endDate).Scan(&report.TotalRevenue, &report.TotalTransaction)
	if err != nil {
		return nil, err
	}

	queryTop := `
		SELECT 
			p.name, 
			SUM(td.quantity) as total_qty
		FROM transaction_details td
		JOIN transactions t ON td.transaction_id = t.id
		JOIN products p ON td.product_id = p.id
		WHERE t.created_at >= $1 AND t.created_at <= $2
		GROUP BY p.name
		ORDER BY total_qty DESC
		LIMIT 1`

	err = repo.db.QueryRow(queryTop, startDate, endDate).Scan(&report.TopProduct.Name, &report.TopProduct.TotalSold)

	if err == sql.ErrNoRows {
		report.TopProduct = models.BestSellingProduct{Name: "-", TotalSold: 0}
	} else if err != nil {
		return nil, err
	}

	return report, nil
}
