package filerepo

import (
	"encoding/json"
	"example.com/expense_tracker/internal/domain/models"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type FileExpenseRepository struct {
	filePath string
	mu       sync.Mutex
}

func NewFileExpenseRepository(filePath string) *FileExpenseRepository {
	repo := &FileExpenseRepository{filePath: filePath}

	return repo
}

func (r *FileExpenseRepository) generateId() string {
	return fmt.Sprintf("exp-%d", time.Now().UnixNano())
}

func (r *FileExpenseRepository) loadExpensesInternal() ([]*models.Expense, error) {
	if _, err := os.Stat(r.filePath); os.IsNotExist(err) {
		dir := filepath.Dir(r.filePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create directory for expense file: %w", err)
		}
		if err := os.WriteFile(r.filePath, []byte("[]"), 0644); err != nil {
			return nil, fmt.Errorf("failed to create initial expense file: %w", err)
		}
		return []*models.Expense{}, nil
	}

	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read expense file: %w", err)
	}

	var expenses []*models.Expense
	if len(data) == 0 {
		return []*models.Expense{}, nil
	}
	if err := json.Unmarshal(data, &expenses); err != nil {
		return nil, fmt.Errorf("failed to unmarshal expenses from file: %w", err)
	}
	return expenses, nil
}

func (r *FileExpenseRepository) saveExpensesInternal(expenses []*models.Expense) error {
	data, err := json.MarshalIndent(expenses, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal expenses: %w", err)
	}

	if err := os.WriteFile(r.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write expenses to file: %w", err)
	}
	return nil
}

func (r *FileExpenseRepository) Add(expense *models.Expense) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if expense.ID == "" {
		expense.ID = r.generateId()
	}

	expenses, err := r.loadExpensesInternal()
	if err != nil {
		return err
	}
	expenses = append(expenses, expense)
	return r.saveExpensesInternal(expenses)
}

func (r *FileExpenseRepository) GetAll() ([]*models.Expense, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.loadExpensesInternal()
}

func (r *FileExpenseRepository) GetById(id string) (*models.Expense, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	expenses, err := r.loadExpensesInternal()
	if err != nil {
		return nil, err
	}

	for _, expense := range expenses {
		if expense.ID == id {
			return expense, nil
		}
	}

	return nil, fmt.Errorf("expense with id %s not found", id)

}

func (r *FileExpenseRepository) RemoveById(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	expenses, err := r.loadExpensesInternal()
	if err != nil {
		return err
	}

	for i, expense := range expenses {
		if expense.ID == id {
			expenses = append(expenses[:i], expenses[i+1:]...)
			break
		}
	}

	return r.saveExpensesInternal(expenses)
}

func (r *FileExpenseRepository) Update(expense *models.Expense) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	expenses, err := r.loadExpensesInternal()
	if err != nil {
		return fmt.Errorf("error getting expenses: %w", err)
	}

	if len(expenses) == 0 {
		return fmt.Errorf("no expenses found")
	}

	for i, exp := range expenses {
		if exp.ID == expense.ID {
			expenses[i] = expense
			return r.saveExpensesInternal(expenses)
		}
	}

	return fmt.Errorf("expense with id %s not found", expense.ID)
}

func (r *FileExpenseRepository) GetByDateRange(fromDate, toDate time.Time) ([]*models.Expense, error) {
	expenses, err := r.GetAll()
	if err != nil {
		return nil, err
	}

	var filtered []*models.Expense

	for _, exp := range expenses {
		if exp.Date.After(fromDate) && exp.Date.Before(toDate) || exp.Date.Equal(fromDate) || exp.Date.Equal(toDate) {
			filtered = append(filtered, exp)
		}
	}

	return filtered, nil
}
