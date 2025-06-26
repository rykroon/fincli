package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/rykroon/ry-cli/internal"
)

var amortizationCmd = &cobra.Command{
	Use:   "amortization",
	Short: "Amortization Calculation.",
	Long:  ``,
	Run:   runAmortizationCmd,
}

type amortizationParams struct {
	Principal     float64
	Rate          float64
	Years         int
	HomePrice     float64
	DownPercent   float64
	PrintSchedule bool
}

var amp amortizationParams

func runAmortizationCmd(cmd *cobra.Command, args []string) {
	p := amp.Principal
	if p < 0 {
		p = amp.HomePrice - (amp.HomePrice * (amp.DownPercent / 100))
	}
	r := amp.Rate / 12 / 100
	n := amp.Years * 12

	monthlyPayment := internal.AmortizedPayment(p, r, n)
	fmt := message.NewPrinter(language.English)
	fmt.Printf("Monthly Payment: $%.2f\n", monthlyPayment)

	if amp.PrintSchedule {
		payments := internal.AmortizationSchedule(p, r, n)

		// header
		fmt.Printf("%-6s %-6s %-12s %-12s %-12s %-12s\n", "Year", "Month", "Payment", "Principal", "Interest", "Balance")
		fmt.Println(strings.Repeat("-", 58))

		for _, payment := range payments {
			// add cumalive interest and equity
			year := payment.Month/12 + 1
			fmt.Printf("%-6d %-6d $%-11.2f $%-11.2f $%-11.2f $%-11.2f\n",
				year, payment.Month, monthlyPayment, payment.PrincipalPaid, payment.InterestPaid, payment.Balance)
		}
	}
}

func init() {
	amortizationCmd.Flags().Float64VarP(&amp.Principal, "principal", "p", -1, "The loan amount borrowed.")
	amortizationCmd.Flags().Float64VarP(&amp.Rate, "rate", "r", -1, "Annual interest rate.")
	amortizationCmd.Flags().IntVarP(&amp.Years, "years", "y", 30, "Loan term in years.")

	// Helper Arguments. Mutually exclusive from the above.
	amortizationCmd.Flags().Float64VarP(&amp.HomePrice, "price", "H", -1, "Home purchase price.")
	amortizationCmd.Flags().Float64VarP(&amp.DownPercent, "down-percent", "d", -1, "Down payment as a percent.")

	amortizationCmd.MarkFlagRequired("rate")
	amortizationCmd.MarkFlagsRequiredTogether("price", "down-percent")
	amortizationCmd.MarkFlagsMutuallyExclusive("principal", "price")
	amortizationCmd.MarkFlagsMutuallyExclusive("principal", "down-percent")

	amortizationCmd.Flags().BoolVarP(&amp.PrintSchedule, "schedule", "s", false, "Print the schedule.")
}
