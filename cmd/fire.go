package cmd

import (
	"fmt"

	"github.com/rykroon/fincli/internal/flagx"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
)

func NewFireCmd() *cobra.Command {
	annualExpenses := decimal.Zero
	safeWithdrawalRate := decimal.NewFromFloat(.04)

	cmd := &cobra.Command{
		Use:   "fire",
		Short: "Calculate your FIRE number.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !annualExpenses.IsPositive() {
				return fmt.Errorf("expenses must be greater than zero")
			}
			if !safeWithdrawalRate.IsPositive() {
				return fmt.Errorf("safe withdrawal rate must be greater than zero")
			}
			fireNumber := annualExpenses.Div(safeWithdrawalRate)
			prt.Printf("%-21s $%13.2v\n", "FIRE Number:", fireNumber)
			prt.Printf(
				"%-22s %13.2v%%\n",
				"Safe Withdrawal Rate:",
				safeWithdrawalRate.Mul(decimal.NewFromInt(100)),
			)
			return nil
		},
	}

	cmd.Flags().VarP(
		flagx.NewDecimalFlag(&annualExpenses), "expenses", "e", "Annual expenses.",
	)

	cmd.Flags().Var(
		flagx.NewPercentFlag(&safeWithdrawalRate),
		"swr",
		"Safe withdrawal rate.",
	)

	cmd.MarkFlagRequired("expenses")

	cmd.Flags().SortFlags = false

	return cmd
}
