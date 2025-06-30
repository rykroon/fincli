package realestate

import (
	"strings"

	"github.com/rykroon/fincli/internal/mortgage"
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
	Amount              float64
	Rate                float64
	Years               int
	ExtraMonthlyPayment float64
	ExtraAnnualPayment  float64
	MonthlySchedule     bool
	AnnualSchedule      bool
}

func (mpf *mortgagePayoffFlags) HasExtraPayment() bool {
	return mpf.ExtraAnnualPayment > 0 || mpf.ExtraMonthlyPayment > 0
}

func (mpf *mortgagePayoffFlags) ExtraPaymentStrategy() mortgage.ExtraPaymentStrategy {
	if mpf.ExtraMonthlyPayment > 0 {
		return mortgage.ExtraMonthlyPaymentStrategy(mpf.ExtraMonthlyPayment)
	} else if mpf.ExtraAnnualPayment > 0 {
		return mortgage.ExtraAnnualPaymentStrategy(mpf.ExtraAnnualPayment, true)
	} else {
		return mortgage.NoExtraPaymentStrategy()
	}
}

var mpf mortgagePayoffFlags

func runMortgageCmd(cmd *cobra.Command, args []string) {
	fmt := message.NewPrinter(language.English)
	i := mpf.Rate / 12 / 100
	n := mpf.Years * 12
	monthlyPayment, payments := mortgage.CalculatePayments(mpf.Amount, i, n, mpf.ExtraPaymentStrategy())
	stats := mortgage.GetPaymentScheduleStats(payments)

	if !mpf.HasExtraPayment() {
		fmt.Printf("Monthly Payment: %.2f\n", monthlyPayment)
	}

	fmt.Printf("Total Amount Paid: $%.2f\n", stats.TotalPayments)
	fmt.Printf("Total Interest Paid: $%.2f\n", stats.TotalInterest)
	fmt.Printf("Pay off in %d years and %d month(s)\n", len(payments)/12, len(payments)%12)
	fmt.Println("")
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
	mortgagePayoffCmd.Flags().Float64Var(&mpf.ExtraMonthlyPayment, "extra-monthly", 0, "Extra monthly payment.")
	mortgagePayoffCmd.Flags().Float64Var(&mpf.ExtraAnnualPayment, "extra-annual", 0, "Extra annual payment.")
	mortgagePayoffCmd.Flags().BoolVar(&mpf.MonthlySchedule, "monthly-schedule", false, "Print the monthly amortization schedule.")
	mortgagePayoffCmd.Flags().BoolVar(&mpf.AnnualSchedule, "annual-schedule", false, "Print the annual amortization schedule.")

	// technically not mutually exclusive, but would need extra logic.
	mortgagePayoffCmd.MarkFlagsMutuallyExclusive("extra-monthly", "extra-annual")
	mortgagePayoffCmd.MarkFlagsMutuallyExclusive("monthly-schedule", "annual-schedule")
}
