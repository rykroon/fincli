package amortization

import (
	"github.com/spf13/cobra"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var paymentCmd = &cobra.Command{
	Use:   "payment",
	Short: "Calculate monthly payment.",
	Long:  ``,
	Run:   runPaymentCmd,
}

var ampp amortizationParams

func runPaymentCmd(cmd *cobra.Command, args []string) {
	fmt := message.NewPrinter(language.English)
	payment := getAmortizedPayment(ampp.Principal(), ampp.MonthlyRate(), ampp.NumPeriods())
	fmt.Printf("Monthly Payment: $%.2f\n", payment)
}

func init() {
	addFlags(paymentCmd, &ampp)
}
