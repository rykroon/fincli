package cmd

import (
	"github.com/rykroon/fincli/internal/cli"
	"github.com/rykroon/fincli/internal/flag"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
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
	fmt := message.NewPrinter(language.English)
	fireNumber := ff.AnnualExpenses.Div(ff.SafeWithdrawlRate)
	fmt.Println("FIRE Number: ", cli.FormatMoney(fireNumber, sep))
}

func init() {
	fireCmd.Flags().VarP(flag.NewDecVal(&ff.AnnualExpenses), "expenses", "e", "Annual expenses.")
	ff.SafeWithdrawlRate = decimal.NewFromFloat(.04)
	fireCmd.Flags().Var(flag.NewPercentVal(&ff.SafeWithdrawlRate), "swr", "Safe withdrawl rate.")

	fireCmd.MarkFlagRequired("expenses")

	fireCmd.Flags().SortFlags = false
	fireCmd.Flags().PrintDefaults()
}
