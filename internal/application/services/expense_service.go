package services

import (
	"example.com/expense_tracker/internal/domain/contracts"
	"example.com/expense_tracker/internal/domain/models"
	"fmt"
	"time"
)

type ExpenseService struct {
	repo contracts.ExpenseRepository
}

func NewExpenseService(repo contracts.ExpenseRepository) *ExpenseService {
	return &ExpenseService{repo: repo}
}

func (s *ExpenseService) AddExpense(description string, amount float64) (*models.Expense, error) {
	expense, err := models.NewExpense(description, amount, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to create new expense: %w", err)
	}

	if err := expense.IsValid(); err != nil {
		return nil, fmt.Errorf("invalid expense data: %w", err)
	}

	if err := s.repo.Add(expense); err != nil {
		return nil, fmt.Errorf("error saving expense to repository: %w", err)
	}

	return expense, nil
}

func (s *ExpenseService) GetById(ID string) (*models.Expense, error) {
	expenses, err := s.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("error getting expenses: %w", err)
	}

	if len(expenses) == 0 {
		return nil, fmt.Errorf("no expenses found")
	}

	for i, exp := range expenses {
		if exp.ID == ID {
			return expenses[i], nil
		}
	}

	return nil, fmt.Errorf("expense with ID %s not found", ID)
}

func (s *ExpenseService) GetAll() ([]*models.Expense, error) {
	expenses, err := s.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("error getting expenses: %w", err)
	}

	return expenses, nil
}

func (s *ExpenseService) RemoveExpense(ID string) error {
	if err := s.repo.RemoveById(ID); err != nil {
		return fmt.Errorf("error removing expense with ID %s: %w", ID, err)
	}

	return nil
}

//func (s *ExpenseService) Delete() {
//
//}
