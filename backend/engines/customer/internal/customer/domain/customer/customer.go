package customer

import (
	"fmt"
	"time"
)

type CustomerID string

type Customer struct {
	ID        CustomerID     `json:"id"`
	Phone     string         `json:"phone"`
	Name      string         `json:"name"`
	Email     string         `json:"email,omitempty"`
	Status    CustomerStatus `json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type CustomerStatus string

const (
	CustomerActive   CustomerStatus = "active"
	CustomerInactive CustomerStatus = "inactive"
	CustomerBlocked  CustomerStatus = "blocked"
)

func NewCustomer(phone string) *Customer {
	now := time.Now()
	return &Customer{
		ID:        CustomerID(id()),
		Phone:     phone,
		Status:    CustomerActive,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func id() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
