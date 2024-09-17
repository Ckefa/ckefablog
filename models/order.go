package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Order represents an order with a package and status information.
type Order struct {
	gorm.Model
	ID        string `json:"id"`
	PackageID int64  `json:"package_id"` // Foreign key for the Package model
	Amount    string `json:"amount"`
	Status    bool   `json:"status"`
	PayId     string `json:"pay_id"`
	StatusMsg string `json:"status_msg"`
	Link1     string `json:"link1"`
	Link2     string `json:"link2"`
}

// NewOrder creates a new order with a generated UUID and a package reference.
func NewOrder(amt string, pid int64) *Order {
	return &Order{
		ID:        uuid.New().String(),
		PackageID: pid,
		Amount:    amt,
		Status:    false,
	}
}
