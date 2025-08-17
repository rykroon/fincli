package cmd

import (
	"github.com/rykroon/fincli/internal/cli"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var fireCmd = &cobra.Command{
	Use:   "fire",
	Short: "Calculate your FIRE number.",
	Run:   runFireNumberCmd,
}

type fireNumberFlags struct {
	AnnualExpenses    decimal.Decimal
	SafeWithdrawlRate decimal.Decimal
}

var fnf fireNumberFlags

func runFireNumberCmd(cmd *cobra.Command, args []string) {
	fmt := message.NewPrinter(language.English)
	fireNumber := fnf.AnnualExpenses.Div(fnf.SafeWithdrawlRate)
	fmt.Println("FIRE Number: ", cli.FormatMoney(fireNumber, sep))
}

func init() {
	fireCmd.Flags().VarP(cli.DecimalValue(&fnf.AnnualExpenses), "expenses", "e", "Annual expenses.")
	fnf.SafeWithdrawlRate = decimal.NewFromFloat(.04)
	fireCmd.Flags().Var(cli.PercentValue(&fnf.SafeWithdrawlRate), "swr", "Safe withdrawl rate.")

	fireCmd.MarkFlagRequired("expenses")

	fireCmd.Flags().SortFlags = false
	fireCmd.Flags().PrintDefaults()
}
