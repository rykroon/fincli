package realestate

import (
	"github.com/spf13/cobra"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var cashNeededCmd = &cobra.Command{
	Use:   "cash-needed",
	Short: "Calculate the amount of cash needed to purchase a home.",
	Run:   runCashNeededCmd,
}

type cashNeededFlags struct {
	Price          float64
	DownPaymentPct float64
	ClosingPct     float64
}

var cnf cashNeededFlags

func runCashNeededCmd(cmd *cobra.Command, args []string) {
	fmt := message.NewPrinter(language.English)
	downPayment := cnf.Price * cnf.DownPaymentPct / 100
	closingCosts := cnf.Price * cnf.ClosingPct / 100
	totalCashNeeded := downPayment + closingCosts

	fmt.Printf("Down payment: $%.2f\n", downPayment)
	fmt.Printf("Closing costs: $%.2f\n", closingCosts)
	fmt.Printf("Total Cash needed: $%.2f\n", totalCashNeeded)
}

func init() {
	cashNeededCmd.Flags().Float64VarP(&cnf.Price, "price", "p", 0, "Home price")
	cashNeededCmd.Flags().Float64VarP(&cnf.DownPaymentPct, "down", "d", 0, "Down payment %")
	cashNeededCmd.Flags().Float64Var(&cnf.ClosingPct, "closing-percent", 0, "Closing cost estimate %")
	// escrow --escrow
	// custom fees --fee name:amount
	cashNeededCmd.MarkFlagRequired("price")
}
