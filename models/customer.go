package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Customer struct {
	gorm.Model
	ID        string  `json:"id"`
	Email     string  `gorm:"unique" json:"email"`
	Fname     string  `json:"fname"`
	Lname     string  `json:"lname"`
	Passwd    string  `json:"passwd"`
	PackageID int64   `json:"package_id"`
	Package   Package `gorm:"foreignKey:PackageID"`
}

func NewCustomer(fname string, lname string, email string, passwd string) *Customer {
	return &Customer{
		ID:        uuid.New().String(),
		Fname:     fname,
		Lname:     lname,
		Email:     email,
		Passwd:    passwd,
		PackageID: 1,
	}
}
