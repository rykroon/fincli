package cmd

import (
	"fmt"
	"math"
	"strings"

	"github.com/spf13/cobra"
)

var amortizeCmd = &cobra.Command{
	Use:   "amortize",
	Short: "Amortization Calculation.",
	Long:  ``,
	Run:   runAmortizeCmd,
}

var p float64           // principal
var homePrice float64   // home price
var downPercent float64 // down payment as percent
var i float64           // monthly interest
var annualRate float64  // annual interest
var n int               // number of payments
var years int           // number of years
var schedule bool

func runAmortizeCmd(cmd *cobra.Command, args []string) {
	if p < 0 {
		p = homePrice - (homePrice * (downPercent / 100))
	}
	if i < 0 {
		i = annualRate / 12
	}
	i = i / 100 // needs to be expressed as a percentage
	if n < 0 {
		n = years * 12
	}

	monthlyPayment := p * ((i * math.Pow(1+i, float64(n))) / (math.Pow(1+i, float64(n)) - 1))
	fmt.Printf("Monthly Payment: $%.2f\n", monthlyPayment)

	balance := p

	if schedule {
		// header
		fmt.Printf("%-6s %-12s %-12s %-12s %-12s\n", "Month", "Payment", "Principal", "Interest", "Balance")
		fmt.Println(strings.Repeat("-", 58))

		for month := 1; month < n; month++ {
			interestPaid := balance * i
			principalPaid := monthlyPayment - interestPaid
			balance -= principalPaid
			fmt.Printf("%-6d $%-11.2f $%-11.2f $%-11.2f $%-11.2f\n",
				month, monthlyPayment, principalPaid, interestPaid, balance)
		}
	}
}

func init() {
	amortizeCmd.Flags().Float64VarP(&p, "principal", "p", -1, "The loan amount borrowed.")
	amortizeCmd.Flags().Float64VarP(&i, "monthly-rate", "i", -1, "The monthly interesst rate.")
	amortizeCmd.Flags().IntVarP(&n, "payments", "n", -1, "Total number of payments.")

	// Helper Arguments. Mutually exclusive from the above.
	amortizeCmd.Flags().Float64VarP(&homePrice, "price", "H", -1, "Home purchase price.")
	amortizeCmd.Flags().Float64VarP(&downPercent, "down-percent", "d", -1, "Down payment as a percent.")
	amortizeCmd.Flags().Float64VarP(&annualRate, "annual-rate", "r", -1, "Annual interest rate.")
	amortizeCmd.Flags().IntVarP(&years, "years", "y", -1, "Loan term in years.")

	amortizeCmd.MarkFlagsRequiredTogether("price", "down-percent")
	amortizeCmd.MarkFlagsMutuallyExclusive("principal", "price")
	amortizeCmd.MarkFlagsMutuallyExclusive("principal", "down-percent")

	amortizeCmd.MarkFlagsMutuallyExclusive("payments", "years")
	amortizeCmd.MarkFlagsOneRequired("payments", "years")

	amortizeCmd.MarkFlagsMutuallyExclusive("monthly-rate", "annual-rate")
	amortizeCmd.MarkFlagsOneRequired("monthly-rate", "annual-rate")

	amortizeCmd.Flags().BoolVarP(&schedule, "schedule", "s", false, "Print the schedule.")
}
