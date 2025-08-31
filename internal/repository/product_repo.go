package repository

import (
	"database/sql"
	"sales-system/internal/models"
	"time"
)

type ProductRepo struct {
	db *sql.DB
}

func NewProductRepo(db *sql.DB) *ProductRepo {
	return &ProductRepo{db: db}
}

func (r *ProductRepo) CreateProduct(p models.Product) (int64, error) {
	res, err := r.db.Exec("INSERT INTO products (date, name, quantity, price) VALUES (?, ?, ?, ?)", p.Date.Format(time.RFC3339), p.Name, p.Quantity, p.Price)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return id, err
}

func (r *ProductRepo) GetProductByID(id int) (*models.Product, error) {
	row := r.db.QueryRow("SELECT id, date, name, quantity, price FROM products WHERE id = ?", id)
	var p models.Product
	var dateStr string
	err := row.Scan(&p.ID, &dateStr, &p.Name, &p.Quantity, &p.Price)
	if err != nil {
		return nil, err
	}
	p.Date, _ = time.Parse(time.RFC3339, dateStr)
	return &p, nil
}

func (r *ProductRepo) GetAllProducts() ([]models.Product, error) {
	rows, err := r.db.Query("SELECT id, date, name, quantity, price FROM products ORDER BY id DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		var dateStr string
		if err := rows.Scan(&p.ID, &dateStr, &p.Name, &p.Quantity, &p.Price); err != nil {
			return nil, err
		}
		p.Date, _ = time.Parse(time.RFC3339, dateStr)
		products = append(products, p)
	}
	return products, nil
}

func (r *ProductRepo) UpdateProduct(p models.Product) error {
	_, err := r.db.Exec("UPDATE products SET date = ?, quantity = ?, price = ? WHERE id = ?", p.Date.Format(time.RFC3339), p.Quantity, p.Price, p.ID)
	return err
}

func (r *ProductRepo) DeleteProduct(id int) error {
	_, err := r.db.Exec("DELETE FROM products WHERE id = ?", id)
	return err
}