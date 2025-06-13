package cmd

import (
	"example.com/expense_tracker/internal/application/services"
	"example.com/expense_tracker/internal/infra/persistence/filerepo"
	"example.com/expense_tracker/pkg/config"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"log"
	"os"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Show list",
	Long:  `Show expense list`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig()
		if err != nil {
			log.Fatal(err)
		}

		expenseRepo := filerepo.NewFileExpenseRepository(cfg.ExpensesFilePath)
		expenseService := services.NewExpenseService(expenseRepo)
		expenses, err := expenseService.GetAll()

		if err != nil {
			fmt.Println(err)
			return
		}

		if len(expenses) == 0 {
			fmt.Println("No expenses found!")
			return
		}

		tableOutput := tablewriter.NewWriter(os.Stdout)
		tableOutput.Header([]string{"№", "ID", "Amount", "Description", "Date"})

		i := 1
		for _, expense := range expenses {
			tableOutput.Append([]string{
				fmt.Sprintf("%d", i),
				expense.ID,
				fmt.Sprintf("%.2f$", expense.Amount),
				expense.Description,
				expense.Date.Format("2006-01-02 15:04:05"),
			})
			i++
		}

		tableOutput.Render()
	},
}
