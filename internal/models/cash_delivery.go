package models

import "time"

type CashDelivery struct {
	ID          int
	Date        time.Time
	Name        string
	Description string
	Amount      float64
}