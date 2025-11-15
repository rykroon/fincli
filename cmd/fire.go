package cmd

import (
	"github.com/rykroon/fincli/internal/flagx"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
)

func NewFireCmd() *cobra.Command {
	var annualExpenses decimal.Decimal
	var safeWithdrawlRate decimal.Decimal

	cmd := &cobra.Command{
		Use:   "fire",
		Short: "Calculate your FIRE number.",
		Run: func(cmd *cobra.Command, args []string) {
			runFireCmd(annualExpenses, safeWithdrawlRate)
		},
	}

	flagx.DecimalVarP(
		cmd.Flags(), &annualExpenses, "expenses", "e", decimal.Zero, "Annual expenses.",
	)

	flagx.PercentVar(
		cmd.Flags(),
		&safeWithdrawlRate,
		"swr",
		decimal.NewFromFloat(.04),
		"Safe withdrawl rate.",
	)

	cmd.MarkFlagRequired("expenses")

	cmd.Flags().SortFlags = false
	cmd.Flags().PrintDefaults()

	return cmd
}

func runFireCmd(annualExpenses, swr decimal.Decimal) {
	fireNumber := annualExpenses.Div(swr)
	prt.Printf("%-20s $%13.2v\n", "FIRE Number:", fireNumber)
	prt.Printf("%-20s %13.2v%%\n", "Safe Withdrawl Rate:", swr.Mul(decimal.NewFromInt(100)))
}
