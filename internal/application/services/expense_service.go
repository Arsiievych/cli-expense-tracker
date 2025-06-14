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

type ExpensesSummary struct {
	Count int
	Sum   float64
}

func NewExpenseService(repo contracts.ExpenseRepository) *ExpenseService {
	return &ExpenseService{repo: repo}
}

func (s *ExpenseService) AddExpense(description string, amount float64, date time.Time) (*models.Expense, error) {
	expense, err := models.NewExpense(description, amount, date)
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
	if ID == "" {
		return nil, fmt.Errorf("invalid expense id")
	}

	return s.repo.GetById(ID)
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

func (s *ExpenseService) UpdateExpense(expense *models.Expense) error {
	if expense == nil {
		return fmt.Errorf("expense is not provided")
	}

	err := s.repo.Update(expense)
	if err != nil {
		return fmt.Errorf("error updating expense with ID %s: %w", expense.ID, err)
	}

	return nil
}

func (s *ExpenseService) GetExpensesByDateRange(from, to time.Time) ([]*models.Expense, error) {
	return s.repo.GetByDateRange(from, to)
}

func (s *ExpenseService) GetExpensesSummary(expenses []*models.Expense) (ExpensesSummary, error) {
	var summary ExpensesSummary = ExpensesSummary{
		Count: len(expenses),
		Sum:   0,
	}

	for _, exp := range expenses {
		summary.Sum += exp.Amount
	}

	return summary, nil
}
