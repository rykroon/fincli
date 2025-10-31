package cmd

import (
	"github.com/rykroon/fincli/internal/flagx"
	"github.com/rykroon/fincli/internal/fmtx"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
)

func NewFireCmd() *cobra.Command {
	var annualExpenses decimal.Decimal
	safeWithdrawlRate := decimal.NewFromFloat(.04)

	cmd := &cobra.Command{
		Use:   "fire",
		Short: "Calculate your FIRE number.",
		Run: func(cmd *cobra.Command, args []string) {
			sep, _ := flagx.GetRune(cmd.PersistentFlags(), "sep")
			prt := fmtx.NewDecimalPrinter(sep)
			runFireCmd(prt, annualExpenses, safeWithdrawlRate)
		},
	}

	flagx.DecimalVarP(cmd.Flags(), &annualExpenses, "expenses", "e", decimal.Zero, "Annual expenses.")
	cmd.Flags().Var(flagx.NewPercentVal(&safeWithdrawlRate), "swr", "Safe withdrawl rate.")

	cmd.MarkFlagRequired("expenses")

	cmd.Flags().SortFlags = false
	cmd.Flags().PrintDefaults()

	return cmd
}

func runFireCmd(prt fmtx.DecimalPrinter, annualExpenses, swr decimal.Decimal) {
	fireNumber := annualExpenses.Div(swr)
	prt.Printf("FIRE Number: $%.2v\n", fireNumber)
	prt.Printf("Safe Withdrawl Rate: %.2v%%\n", swr.Mul(decimal.NewFromInt(100)))
}
