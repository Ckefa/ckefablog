package models

import (
	"github.com/google/uuid"
)

type Message struct {
	ID         string `json:"id"`          // Set ID as varchar and primary key
	CustomerID string `json:"customer_id"` // Link to CustomerID
	Body       string `json:"body"`
	Role       string `json:"role"`
}

func NewMessage(cid string, msg string, role string) *Message {
	return &Message{
		ID:         uuid.New().String(),
		CustomerID: cid,
		Body:       msg,
		Role:       role,
	}
}
