package repository

import (
	"database/sql"
	"sales-system/internal/models"
	"time"
)

type CashDeliveryRepo struct {
	db *sql.DB
}

func NewCashDeliveryRepo(db *sql.DB) *CashDeliveryRepo {
	return &CashDeliveryRepo{db: db}
}

func (r *CashDeliveryRepo) CreateCashDelivery(cd models.CashDelivery) (int64, error) {
	res, err := r.db.Exec("INSERT INTO cash_deliveries (date, name, description, amount) VALUES (?, ?, ?, ?)", cd.Date.Format(time.RFC3339), cd.Name, cd.Description, cd.Amount)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return id, err
}

func (r *CashDeliveryRepo) GetCashDeliveriesByDateRange(start, end time.Time) ([]models.CashDelivery, error) {
	rows, err := r.db.Query("SELECT id, date, name, description, amount FROM cash_deliveries WHERE date BETWEEN ? AND ?", start.Format(time.RFC3339), end.Format(time.RFC3339))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deliveries []models.CashDelivery
	for rows.Next() {
		var cd models.CashDelivery
		var dateStr string
		if err := rows.Scan(&cd.ID, &dateStr, &cd.Name, &cd.Description, &cd.Amount); err != nil {
			return nil, err
		}
		cd.Date, _ = time.Parse(time.RFC3339, dateStr)
		deliveries = append(deliveries, cd)
	}
	return deliveries, nil
}