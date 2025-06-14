package cmd

import (
	"example.com/expense_tracker/internal/application/services"
	"example.com/expense_tracker/internal/infra/persistence/filerepo"
	"example.com/expense_tracker/pkg/config"
	"fmt"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
	"github.com/spf13/cobra"
	"log"
	"os"
)

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

		colorCfg := renderer.ColorizedConfig{
			Header: renderer.Tint{
				FG: renderer.Colors{color.FgGreen, color.Bold},
			},
			Column: renderer.Tint{
				FG: renderer.Colors{color.FgCyan},
				Columns: []renderer.Tint{
					{FG: renderer.Colors{color.FgWhite}},
					{},
					{FG: renderer.Colors{color.FgHiRed}},
				},
			},
			Border:    renderer.Tint{FG: renderer.Colors{color.FgWhite}},
			Separator: renderer.Tint{FG: renderer.Colors{color.FgWhite}},
		}

		tableOutput := tablewriter.NewTable(os.Stdout,
			tablewriter.WithRenderer(renderer.NewColorized(colorCfg)),
			tablewriter.WithConfig(tablewriter.Config{
				Row: tw.CellConfig{
					Formatting: tw.CellFormatting{
						AutoWrap: tw.WrapNormal,
					},
					Alignment: tw.CellAlignment{
						Global: tw.AlignLeft,
						PerColumn: []tw.Align{
							tw.AlignLeft,
							tw.AlignCenter,
							tw.AlignRight,
							tw.AlignLeft,
							tw.AlignCenter,
						},
					},
					ColMaxWidths: tw.CellWidth{Global: 40},
				},
			}),
		)

		tableOutput.Header([]string{"№", "ID", "Amount", "Description", "Date"})

		i := 1
		for _, expense := range expenses {
			tableOutput.Append([]string{
				fmt.Sprintf("%d", i),
				expense.ID,
				fmt.Sprintf("%.2f$", expense.Amount),
				expense.Description,
				expense.Date.Format("Mon, Jan 2, 2006 15:04:05"),
			})
			i++
		}

		tableOutput.Render()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
