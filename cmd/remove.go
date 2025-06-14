package cmd

import (
	"example.com/expense_tracker/internal/application/services"
	"example.com/expense_tracker/internal/infra/persistence/filerepo"
	"example.com/expense_tracker/pkg/config"
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove expense",
	Long:  `Remove expense from the list`,
	Run: func(cmd *cobra.Command, args []string) {
		id, err := cmd.Flags().GetString("id")
		if err != nil {
			log.Fatal(err)
		}

		config, err := config.LoadConfig()
		if err != nil {
			log.Fatal(err)
		}
		expenseRepo := filerepo.NewFileExpenseRepository(config.ExpensesFilePath)
		expenseService := services.NewExpenseService(expenseRepo)

		if err := expenseService.RemoveExpense(id); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Expense %s has removed!\n", id)
	},
}

func init() {
	removeCmd.Flags().String("id", "id", "expense ID to remove")
	removeCmd.MarkFlagRequired("id")
	rootCmd.AddCommand(removeCmd)
}
