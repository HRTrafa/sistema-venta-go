package models

import "time"

const (
	StatusPaid    = "Pagado"
	StatusPending = "Pendiente"
)

type Sale struct {
	ID        int
	Date      time.Time
	Client    string
	ProductID int
	Quantity  int
	Price     float64
	Total     float64
	Status    string
}