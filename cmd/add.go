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
	addCmd.Flags().StringP("description", "d", "Other", "Expense description")
	addCmd.Flags().Float64P("amount", "a", 0, "Expense amount (e.g., 9.99)")

	addCmd.MarkFlagRequired("amount")

	rootCmd.AddCommand(addCmd)
}
