package contracts

import "example.com/expense_tracker/internal/domain/models"

type ExpenseRepository interface {
	Add(expense *models.Expense) error
	GetAll() ([]*models.Expense, error)
	GetById(id string) (*models.Expense, error)
	RemoveById(id string) error
	Update(expense *models.Expense) error
}
