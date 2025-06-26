package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/rykroon/ry-cli/internal"
)

var amortizeCmd = &cobra.Command{
	Use:   "amortize",
	Short: "Amortization Calculation.",
	Long:  ``,
	Run:   runAmortizeCmd,
}

var p float64
var rate float64
var years int

var homePrice float64
var downPercent float64

var printSchedule bool

func runAmortizeCmd(cmd *cobra.Command, args []string) {
	if p < 0 {
		p = homePrice - (homePrice * (downPercent / 100))
	}
	r := rate / 12 / 100
	n := years * 12

	monthlyPayment := internal.AmortizeLoan(p, r, n)
	fmt := message.NewPrinter(language.English)
	fmt.Printf("Monthly Payment: $%.2f\n", monthlyPayment)

	if printSchedule {
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
	amortizeCmd.Flags().Float64VarP(&p, "principal", "p", -1, "The loan amount borrowed.")
	amortizeCmd.Flags().Float64VarP(&rate, "rate", "r", -1, "Annual interest rate.")
	amortizeCmd.Flags().IntVarP(&years, "years", "y", 30, "Loan term in years.")

	// Helper Arguments. Mutually exclusive from the above.
	amortizeCmd.Flags().Float64VarP(&homePrice, "price", "H", -1, "Home purchase price.")
	amortizeCmd.Flags().Float64VarP(&downPercent, "down-percent", "d", -1, "Down payment as a percent.")

	amortizeCmd.MarkFlagsRequiredTogether("price", "down-percent")
	amortizeCmd.MarkFlagsMutuallyExclusive("principal", "price")
	amortizeCmd.MarkFlagsMutuallyExclusive("principal", "down-percent")

	amortizeCmd.Flags().BoolVarP(&printSchedule, "schedule", "s", false, "Print the schedule.")
}
