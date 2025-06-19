package cmd

import (
	"example.com/expense_tracker/internal/application/services"
	"example.com/expense_tracker/internal/infra/persistence/filerepo"
	"example.com/expense_tracker/pkg/config"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"time"
)

const (
	DateTimeTFormat      = "2006-01-02T15:04:05"
	DateTimeFormat       = "2006-01-02 15:04:05"
	DateTimeFormatShort  = "2006-01-02 15:04"
	DateTimeTFormatShort = "2006-01-02 15:04"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new expense",
	Long:  `Adds a new expense to your expense tracker.`,
	Run: func(cmd *cobra.Command, args []string) {
		description, err := cmd.Flags().GetString("description")
		if err != nil {
			log.Fatal(err)
		}

		amount, err := cmd.Flags().GetFloat64("amount")
		if err != nil {
			log.Fatal(err)
		}

		date := time.Now()
		if cmd.Flags().Changed("time") {
			dateStr, _ := cmd.Flags().GetString("time")
			date, err = dateStrToDate(dateStr)
			if err != nil {
				log.Fatalf("❌ Invalid date format: %v", err)
			}
		}

		cfg, err := config.LoadConfig()
		if err != nil {
			log.Fatal(err)
		}

		expenseRepo := filerepo.NewFileExpenseRepository(cfg.ExpensesFilePath)
		expenseService := services.NewExpenseService(expenseRepo)

		addedExpense, err := expenseService.AddExpense(description, amount, date)
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Failed to add expense: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✅ Successfully added expense: %s\n", addedExpense.ID)
	},
}

func init() {
	addCmd.Flags().StringP("description", "d", "Other", "Expense description")
	addCmd.Flags().Float64P("amount", "a", 0, "Expense amount (e.g., 9.99)")
	addCmd.Flags().StringP("time", "t", "", "Date and time of the expense (e.g., 2025-06-13T15:04)")

	addCmd.MarkFlagRequired("amount")

	rootCmd.AddCommand(addCmd)
}

func dateStrToDate(datetimeStr string) (time.Time, error) {
	if datetimeStr == "" {
		return time.Time{}, fmt.Errorf("date string is empty")
	}

	formats := []string{
		DateTimeFormatShort,
		DateTimeTFormatShort,
		DateTimeFormat,
		DateTimeTFormat,
	}

	for _, format := range formats {
		date, err := time.Parse(format, datetimeStr)
		if err == nil {
			return date, nil
		}
	}

	return time.Time{}, fmt.Errorf("❌ Invalid date format: %v", datetimeStr)
}
