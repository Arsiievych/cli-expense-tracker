package models

import (
	"fmt"
	"time"
)

type Expense struct {
	ID          string    `json:"id"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	Date        time.Time `json:"date"`
}

func NewExpense(description string, amount float64, date time.Time) (*Expense, error) {
	if description == "" {
		return nil, fmt.Errorf("description cannot be empty")
	}

	if amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than zero")
	}

	return &Expense{
		Description: description,
		Amount:      amount,
		Date:        date,
	}, nil
}

func (e *Expense) String() string {
	return fmt.Sprintf("ID: %s, Description: '%s', Amount: $%.2f, Date: %s",
		e.ID, e.Description, e.Amount, e.Date.Format("2006-01-02 15:04:05"))
}

func (e *Expense) IsValid() error {
	return nil
}
