package realestate

import (
	"strings"

	"github.com/rykroon/ry-cli/internal/mortgage"
	"github.com/spf13/cobra"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var mortgagePayoffCmd = &cobra.Command{
	Use:   "mortgage-payoff",
	Short: "Calculate mortgage costs.",
	Run:   runMortgageCmd,
}

type mortgagePayoffFlags struct {
	Amount          float64
	Rate            float64
	Years           int
	ExtraPayment    float64
	MonthlySchedule bool
	AnnualSchedule  bool
}

var mpf mortgagePayoffFlags

func runMortgageCmd(cmd *cobra.Command, args []string) {
	fmt := message.NewPrinter(language.English)
	i := mpf.Rate / 12 / 100
	n := mpf.Years * 12
	// eps := mortgage.NewMonthlyPaymentStrategy(mpf.ExtraPayment)
	monthlyPayment, payments := mortgage.CalculatePayments(mpf.Amount, i, n, mortgage.NewExtraMonthlyPaymentStrategy(749.08))
	monthlyPayment += mpf.ExtraPayment
	stats := mortgage.GetPaymentScheduleStats(payments)
	fmt.Printf("Monthly Payment: %.2f\n", stats.AverageMonthlyPayment)
	fmt.Printf("Total Amount Paid: $%.2f\n", stats.TotalPayments)
	fmt.Printf("Total Interest Paid: $%.2f\n", stats.TotalInterest)
	fmt.Printf("Pay off in %d years and %d month(s)\n", len(payments)/12, len(payments)%12)

	if mpf.MonthlySchedule {
		// print headers
		fmt.Printf(
			"%-6s %-12s %-12s %-12s\n",
			"Month",
			"Principal",
			"Interest",
			"Balance",
		)
		fmt.Println(strings.Repeat("-", 55))
		for _, payment := range payments {
			fmt.Printf("%-6d $%-11.2f $%-11.2f $%-11.2f\n",
				payment.Period,
				payment.PrincipalPaid,
				payment.InterestPaid,
				payment.Balance,
			)

			if payment.Period%12 == 0 {
				fmt.Printf("\t--- End of Year %d ---\n", payment.Period/12)
			}
		}
	} else if mpf.AnnualSchedule {
		// print headers
		fmt.Printf(
			"%-6s %-12s %-12s %-12s\n",
			"Year",
			"Principal",
			"Interest",
			"Balance",
		)
		fmt.Println(strings.Repeat("-", 55))
		annualPrincipalPaid := float64(0)
		annualInterestPaid := float64(0)

		for _, payment := range payments {
			annualPrincipalPaid += payment.PrincipalPaid
			annualInterestPaid += payment.InterestPaid

			if payment.Period%12 == 0 {
				fmt.Printf("%-6d $%-11.2f $%-11.2f $%-11.2f\n",
					payment.Period/12,
					annualPrincipalPaid,
					annualInterestPaid,
					payment.Balance,
				)
				annualPrincipalPaid = 0
				annualInterestPaid = 0
			}
		}

	}
}

func init() {
	mortgagePayoffCmd.Flags().Float64VarP(&mpf.Amount, "amount", "a", 0, "The loan amount borrowed.")
	mortgagePayoffCmd.Flags().Float64VarP(&mpf.Rate, "rate", "r", 0, "Annual interest rate.")
	mortgagePayoffCmd.Flags().IntVarP(&mpf.Years, "years", "y", 30, "Loan term in years.")

	mortgagePayoffCmd.MarkFlagRequired("amount")
	mortgagePayoffCmd.MarkFlagRequired("rate")

	// optional flags
	mortgagePayoffCmd.Flags().Float64Var(&mpf.ExtraPayment, "extra", 0, "Extra payment per month.")
	mortgagePayoffCmd.Flags().BoolVar(&mpf.MonthlySchedule, "monthly-schedule", false, "Print the monthly amortization schedule.")
	mortgagePayoffCmd.Flags().BoolVar(&mpf.AnnualSchedule, "annual-schedule", false, "Print the annual amortization schedule.")

	mortgagePayoffCmd.MarkFlagsMutuallyExclusive("monthly-schedule", "annual-schedule")
}
