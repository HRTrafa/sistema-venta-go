package repository

import (
	"database/sql"
	"sales-system/internal/models"
	"time"
)

type SaleRepo struct {
	db *sql.DB
}

func NewSaleRepo(db *sql.DB) *SaleRepo {
	return &SaleRepo{db: db}
}

func (r *SaleRepo) CreateSale(s models.Sale) (int64, error) {
	res, err := r.db.Exec("INSERT INTO sales (date, client, product_id, quantity, price, total, status) VALUES (?, ?, ?, ?, ?, ?, ?)", s.Date.Format(time.RFC3339), s.Client, s.ProductID, s.Quantity, s.Price, s.Total, s.Status)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return id, err
}

func (r *SaleRepo) GetSaleByID(id int) (*models.Sale, error) {
	row := r.db.QueryRow("SELECT id, date, client, product_id, quantity, price, total, status FROM sales WHERE id = ?", id)
	var s models.Sale
	var dateStr string
	err := row.Scan(&s.ID, &dateStr, &s.Client, &s.ProductID, &s.Quantity, &s.Price, &s.Total, &s.Status)
	if err != nil {
		return nil, err
	}
	s.Date, _ = time.Parse(time.RFC3339, dateStr)
	return &s, nil
}

func (r *SaleRepo) GetAllSales() ([]models.Sale, error) {
	rows, err := r.db.Query("SELECT id, date, client, product_id, quantity, price, total, status FROM sales ORDER BY id DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sales []models.Sale
	for rows.Next() {
		var s models.Sale
		var dateStr string
		if err := rows.Scan(&s.ID, &dateStr, &s.Client, &s.ProductID, &s.Quantity, &s.Price, &s.Total, &s.Status); err != nil {
			return nil, err
		}
		s.Date, _ = time.Parse(time.RFC3339, dateStr)
		sales = append(sales, s)
	}
	return sales, nil
}

func (r *SaleRepo) UpdateSale(s models.Sale) error {
	_, err := r.db.Exec("UPDATE sales SET date = ?, client = ?, quantity = ?, price = ?, total = ?, status = ? WHERE id = ?", s.Date.Format(time.RFC3339), s.Client, s.Quantity, s.Price, s.Total, s.Status, s.ID)
	return err
}

func (r *SaleRepo) DeleteSale(id int) error {
	_, err := r.db.Exec("DELETE FROM sales WHERE id = ?", id)
	return err
}

func (r *SaleRepo) GetSalesByDateRange(start, end time.Time) ([]models.Sale, error) {
	rows, err := r.db.Query("SELECT id, date, client, product_id, quantity, price, total, status FROM sales WHERE date BETWEEN ? AND ?", start.Format(time.RFC3339), end.Format(time.RFC3339))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sales []models.Sale
	for rows.Next() {
		var s models.Sale
		var dateStr string
		if err := rows.Scan(&s.ID, &dateStr, &s.Client, &s.ProductID, &s.Quantity, &s.Price, &s.Total, &s.Status); err != nil {
			return nil, err
		}
		s.Date, _ = time.Parse(time.RFC3339, dateStr)
		sales = append(sales, s)
	}
	return sales, nil
}