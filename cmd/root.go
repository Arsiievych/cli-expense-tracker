package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "expense-tracker",
	Short: "Track your expenses easily from the command line",
	Long:  `CLI Expense Tracker is a simple and fast command-line application for recording and managing your daily expenses.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Expense Tracker CLI - use --help to see available commands")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
