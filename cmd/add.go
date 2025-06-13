package cmd

import (
	"example.com/expense_tracker/internal/application/services"
	"example.com/expense_tracker/internal/infra/persistence/filerepo"
	"example.com/expense_tracker/pkg/config"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var description string
var amount float64

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new expense",
	Long:  `Adds a new expense to your expense tracker.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig()
		if err != nil {
			log.Fatal(err)
		}

		expenseRepo := filerepo.NewFileExpenseRepository(cfg.ExpensesFilePath)
		expenseService := services.NewExpenseService(expenseRepo)

		addedExpense, err := expenseService.AddExpense(description, amount)
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Failed to add expense: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✅ Successfully added expense: %s\n", addedExpense.ID)
	},
}

func init() {
	addCmd.Flags().StringVarP(&description, "desc", "d", "", "Expense description")
	addCmd.Flags().Float64VarP(&amount, "amount", "a", 0, "Expense amount (e.g., 9.99)")

	addCmd.MarkFlagRequired("desc")
	addCmd.MarkFlagRequired("amount")

	rootCmd.AddCommand(addCmd)
}
