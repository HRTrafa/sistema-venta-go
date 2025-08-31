package models

import "time"

type Product struct {
	ID       int
	Date     time.Time
	Name     string
	Quantity int
	Price    float64
}