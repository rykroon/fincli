package investing

import (
	"github.com/spf13/cobra"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var fireNumberCmd = &cobra.Command{
	Use:   "fire-number",
	Short: "Calculate your FIRE number.",
	Run:   runFireNumberCmd,
}

type fireNumberFlags struct {
	AnnualExpenses    float64
	SafeWithdrawlRate float64
}

var fnf fireNumberFlags

func runFireNumberCmd(cmd *cobra.Command, args []string) {
	fmt := message.NewPrinter(language.English)
	fireNumber := fnf.AnnualExpenses / (fnf.SafeWithdrawlRate / 100)
	fmt.Printf("FIRE Number: $%.2f\n", fireNumber)
}

func init() {
	fireNumberCmd.Flags().Float64VarP(&fnf.AnnualExpenses, "expenses", "e", 0, "Annual expenses.")
	fireNumberCmd.Flags().Float64Var(&fnf.SafeWithdrawlRate, "swr", 4, "Safe withdrawl rate.")

	fireNumberCmd.MarkFlagRequired("expenses")
}
