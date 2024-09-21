package models

import (
	"github.com/google/uuid"
)

type Customer struct {
	ID        string    `gorm:"type:varchar(50);primaryKey" json:"id"` // Use string with primary key
	Email     string    `gorm:"unique" json:"email"`
	Fname     string    `json:"fname"`
	Lname     string    `json:"lname"`
	Passwd    string    `json:"passwd"`
	PackageID int64     `json:"package_id"`
	Package   Package   `gorm:"foreignKey:PackageID"`
	Messages  []Message `gorm:"foreignKey:CustomerID" json:"messages"` // Foreign key properly linked
}

func NewCustomer(fname string, lname string, email string, passwd string) *Customer {
	return &Customer{
		ID:        uuid.New().String(),
		Fname:     fname,
		Lname:     lname,
		Email:     email,
		Passwd:    passwd,
		PackageID: 1,
		Messages:  []Message{},
	}
}
