package cmd

import (
	"example.com/expense_tracker/internal/application/services"
	"example.com/expense_tracker/internal/infra/persistence/filerepo"
	"example.com/expense_tracker/pkg/config"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"time"
)

var summaryCmd = &cobra.Command{
	Use:   "summary",
	Short: "Get summary of expenses",
	Long:  `Get summary of expenses`,
	Run: func(cmd *cobra.Command, args []string) {

		cfg, err := config.LoadConfig()
		if err != nil {
			log.Fatal(err)
		}

		expenseRepo := filerepo.NewFileExpenseRepository(cfg.ExpensesFilePath)
		expenseService := services.NewExpenseService(expenseRepo)

		var from, to time.Time

		if cmd.Flags().Changed("month") {
			yearAndMonthStr, err := cmd.Flags().GetString("month")
			if err != nil {
				log.Fatal(err)
			}

			yearAndMonth, err := time.Parse("2006-01", yearAndMonthStr)
			if err != nil {
				log.Fatal(err)
			}

			from = time.Date(yearAndMonth.Year(), yearAndMonth.Month(), 1, 0, 0, 0, 0, time.UTC)
			to = time.Date(yearAndMonth.Year(), yearAndMonth.Month()+1, 0, 23, 59, 59, int(time.Second-time.Nanosecond), time.UTC)
		} else {
			fromStr, err := cmd.Flags().GetString("from")
			if err != nil {
				log.Fatal(err)
			}
			toStr, err := cmd.Flags().GetString("to")
			if err != nil {
				log.Fatal(err)
			}

			if fromStr != "" {
				from, err = dateStrToDate(fromStr)
				if err != nil {
					log.Fatalf("invalid from date: %v", err)
				}
			} else {
				from = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
			}

			if toStr != "" {
				to, err = dateStrToDate(toStr)
				if err != nil {
					log.Fatalf("invalid to date: %v", err)
				}
			} else {
				to = time.Now()
			}
		}

		expenses, err := expenseService.GetExpensesByDateRange(from, to)
		if err != nil {
			log.Fatal(err)
		}

		summary, err := expenseService.GetExpensesSummary(expenses)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("💵 From %s to %s: %d transactions, total amount spent: %.2f$.\n", from.Format("2006-01-02"), to.Format("2006-01-02"), summary.Count, summary.Sum)
	},
}

func init() {
	summaryCmd.Flags().StringP("month", "m", "2025-01", "Summary for a specific month (format: YYYY-MM, e.g., 2025-06)")
	summaryCmd.Flags().StringP("from", "f", "", "Start date for the summary period (format: YYYY-MM, e.g., 2025-06-01 12:00:00)")
	summaryCmd.Flags().StringP("to", "t", "", "End date for the summary period (format: YYYY-MM, e.g., 2025-06-30 14:00:00)")

	rootCmd.AddCommand(summaryCmd)
}
