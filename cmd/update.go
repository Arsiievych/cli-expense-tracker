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

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an expense",
	Long:  `Update an expense from the list.`,
	Run: func(cmd *cobra.Command, args []string) {
		id, err := cmd.Flags().GetString("id")
		if err != nil {
			log.Fatal(err)
		}

		cfg, err := config.LoadConfig()
		if err != nil {
			log.Fatal(err)
		}

		expenseRepo := filerepo.NewFileExpenseRepository(cfg.ExpensesFilePath)
		expenseService := services.NewExpenseService(expenseRepo)

		expense, err := expenseService.GetById(id)
		if err != nil {
			log.Fatal(err)
		}

		if cmd.Flags().Changed("description") {
			description, err := cmd.Flags().GetString("description")
			if err != nil {
				log.Fatal(err)
			}
			expense.Description = description
		}

		if cmd.Flags().Changed("amount") {
			amount, err := cmd.Flags().GetFloat64("amount")
			if err != nil {
				log.Fatal(err)
			}
			expense.Amount = amount
		}

		if cmd.Flags().Changed("time") {
			dateStr, _ := cmd.Flags().GetString("time")
			date, err := dateStrToDate(dateStr)
			if err != nil {
				log.Fatalf("❌ Invalid date format: %v", err)
			}
			expense.Date = date
		}

		if err := expense.IsValid(); err != nil {
			log.Fatal(err)
		}

		err = expenseService.UpdateExpense(expense)
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Failed to update expense: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✅ Successfully Updated expense: %s\n", expense.ID)
	},
}

func init() {
	updateCmd.Flags().String("id", "", "Expense ID")
	updateCmd.Flags().StringP("description", "d", "", "Expense description")
	updateCmd.Flags().Float64P("amount", "a", 0, "Expense amount (e.g., 9.99)")
	updateCmd.Flags().StringP("time", "t", "", "Expense time")

	updateCmd.MarkFlagRequired("id")

	rootCmd.AddCommand(updateCmd)
}
