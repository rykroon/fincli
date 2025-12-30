package cmd

import (
	"github.com/rykroon/fincli/internal/flagx"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
)

func NewFireCmd() *cobra.Command {
	annualExpenses := decimal.Zero
	safeWithdrawlRate := decimal.NewFromFloat(.04)

	cmd := &cobra.Command{
		Use:   "fire",
		Short: "Calculate your FIRE number.",
		Run: func(cmd *cobra.Command, args []string) {
			fireNumber := annualExpenses.Div(safeWithdrawlRate)
			prt.Printf("%-20s $%13.2v\n", "FIRE Number:", fireNumber)
			prt.Printf(
				"%-20s %13.2v%%\n",
				"Safe Withdrawl Rate:",
				safeWithdrawlRate.Mul(decimal.NewFromInt(100)),
			)
		},
	}

	cmd.Flags().VarP(
		flagx.NewDecimalFlag(&annualExpenses), "expenses", "e", "Annual expenses.",
	)

	cmd.Flags().Var(
		flagx.NewPercentFlag(&safeWithdrawlRate),
		"swr",
		"Safe withdrawl rate.",
	)

	cmd.MarkFlagRequired("expenses")

	cmd.Flags().SortFlags = false
	cmd.Flags().PrintDefaults()

	return cmd
}
