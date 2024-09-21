package models

import (
	"fmt"
	"log"
	"strconv"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Order represents an order with a package and status information.
type Order struct {
	gorm.Model
	ID          string `json:"id"`
	PackageID   int64  `json:"package_id"` // Foreign key for the Package model
	PackageName string `json:"package_name"`
	Amount      string `json:"amount"`
	Status      bool   `json:"status"`
	PayId       string `json:"pay_id"`
	StatusMsg   string `json:"status_msg"`
	Link1       string `json:"link1"`
	Link2       string `json:"link2"`
}

// NewOrder creates a new order with a generated UUID and a package reference.
func NewOrder(amt string, pid int64) *Order {
	amount, err := strconv.ParseFloat(amt, 64)
	if err != nil {
		log.Println("<< func: NewOrder  - Error converting amount", err)
	}
	return &Order{
		ID:          uuid.New().String(),
		PackageID:   pid,
		PackageName: Packages[pid-1].Name,
		Amount:      fmt.Sprintf("%.2f", amount),
		Status:      false,
	}
}
