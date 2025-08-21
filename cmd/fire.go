package cmd

import (
	"github.com/rykroon/fincli/internal/flagx"
	"github.com/rykroon/fincli/internal/fmtx"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
)

var fireCmd = &cobra.Command{
	Use:   "fire",
	Short: "Calculate your FIRE number.",
	Run:   runFireCmd,
}

type fireFlags struct {
	AnnualExpenses    decimal.Decimal
	SafeWithdrawlRate decimal.Decimal
}

var ff fireFlags

func runFireCmd(cmd *cobra.Command, args []string) {
	prt := fmtx.NewDecimalPrinter(sep)
	fireNumber := ff.AnnualExpenses.Div(ff.SafeWithdrawlRate)
	prt.Printf("FIRE Number: $%.2v", fireNumber)
}

func init() {
	fireCmd.Flags().VarP(flagx.NewDecVal(&ff.AnnualExpenses), "expenses", "e", "Annual expenses.")
	ff.SafeWithdrawlRate = decimal.NewFromFloat(.04)
	fireCmd.Flags().Var(flagx.NewPercentVal(&ff.SafeWithdrawlRate), "swr", "Safe withdrawl rate.")

	fireCmd.MarkFlagRequired("expenses")

	fireCmd.Flags().SortFlags = false
	fireCmd.Flags().PrintDefaults()
}
