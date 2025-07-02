package realestate

import (
	"strings"

	"github.com/rykroon/fincli/internal/mortgage"
	"github.com/spf13/cobra"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var mortgageCmd = &cobra.Command{
	Use:   "mortgage",
	Short: "Calculate mortgage costs.",
	Run:   runMortgageCmd,
}

type mortgageFlags struct {
	Amount              float64
	Rate                float64
	Years               uint16
	ExtraMonthlyPayment float64
	ExtraAnnualPayment  float64
	MonthlySchedule     bool
	AnnualSchedule      bool
}

func (mf *mortgageFlags) HasExtraPayment() bool {
	return mf.ExtraAnnualPayment > 0 || mf.ExtraMonthlyPayment > 0
}

func (mf *mortgageFlags) ExtraPaymentStrategy() mortgage.ExtraPaymentStrategy {
	// return mortgage.PrincipalMatchInterest()
	if mf.ExtraMonthlyPayment > 0 && mf.ExtraAnnualPayment > 0 {
		return mortgage.ExtraMonthlyAndAnnualPayment(mf.ExtraMonthlyPayment, mf.ExtraAnnualPayment)
	} else if mf.ExtraMonthlyPayment > 0 {
		return mortgage.ExtraMonthlyPayment(mf.ExtraMonthlyPayment)
	} else if mf.ExtraAnnualPayment > 0 {
		return mortgage.ExtraAnnualPayment(mf.ExtraAnnualPayment)
	} else {
		return mortgage.NoExtraPayment()
	}
}

var mf mortgageFlags

func runMortgageCmd(cmd *cobra.Command, args []string) {
	fmt := message.NewPrinter(language.English)
	i := mf.Rate / 12 / 100
	n := int(mf.Years) * 12
	monthlyPayment, payments := mortgage.CalculatePayments(mf.Amount, i, n, mf.ExtraPaymentStrategy())
	stats := mortgage.GetPaymentScheduleStats(payments)

	if !mf.HasExtraPayment() {
		fmt.Printf("Monthly Payment: %.2f\n", monthlyPayment)
	}

	fmt.Printf("Total Amount Paid: $%.2f\n", stats.TotalPayments)
	fmt.Printf("Total Interest Paid: $%.2f\n", stats.TotalInterest)
	fmt.Printf("Pay off in %d years and %d month(s)\n", len(payments)/12, len(payments)%12)
	fmt.Println("")

	if mf.MonthlySchedule {
		printMonthlySchedule(payments)

	} else if mf.AnnualSchedule {
		printAnnualSchedule(payments)
	}
}

func printMonthlySchedule(payments []mortgage.Payment) {
	fmt := message.NewPrinter(language.English)
	fmt.Printf(
		"%-6s %-12s %-12s %-12s %-12s %-12s %-12s\n",
		"Month",
		"Principal",
		"Extra Principal",
		"Total Principal",
		"Interest",
		"Total Payment",
		"Balance",
	)
	fmt.Println(strings.Repeat("-", 89))
	for _, payment := range payments {
		fmt.Printf(
			"%-6d $%-11.2f $%-14.2f $%-14.2f $%-11.2f $%-12.2f $%-11.2f\n",
			payment.Period,
			payment.Principal,
			payment.ExtraPrincipal,
			payment.TotalPrincipal(),
			payment.Interest,
			payment.Payment(),
			payment.Balance,
		)

		if payment.Period%12 == 0 {
			fmt.Printf("\t--- End of Year %d ---\n", payment.Period/12)
		}
	}
}

func printAnnualSchedule(payments []mortgage.Payment) {
	fmt := message.NewPrinter(language.English)
	fmt.Printf(
		"%-6s %-12s %-12s %-12s %-12s %-12s %-12s\n",
		"Year",
		"Principal",
		"Extra Principal",
		"Total Principal",
		"Interest",
		"Total Payment",
		"Balance",
	)
	fmt.Println(strings.Repeat("-", 89))
	var annualPrincipal float64 = 0
	var annualExtraPrincipal float64 = 0
	var annualTotalPrincipal float64 = 0
	var annualInterest float64 = 0
	var annualPayments float64 = 0

	for _, payment := range payments {
		annualPrincipal += payment.Principal
		annualExtraPrincipal += payment.ExtraPrincipal
		annualTotalPrincipal += payment.TotalPrincipal()
		annualInterest += payment.Interest
		annualPayments += payment.Payment()

		if payment.Period%12 == 0 {
			fmt.Printf(
				"%-6d $%-11.2f $%-14.2f $%-14.2f $%-11.2f $%-12.2f $%-11.2f\n",
				payment.Period/12,
				annualPrincipal,
				annualExtraPrincipal,
				annualTotalPrincipal,
				annualInterest,
				annualPayments,
				payment.Balance,
			)
			annualPrincipal = 0
			annualExtraPrincipal = 0
			annualTotalPrincipal = 0
			annualInterest = 0
			annualPayments = 0

		}
	}
}

func init() {
	mortgageCmd.Flags().Float64VarP(&mf.Amount, "amount", "a", 0, "The loan amount borrowed.")
	mortgageCmd.Flags().Float64VarP(&mf.Rate, "rate", "r", 0, "Annual interest rate.")
	mortgageCmd.Flags().Uint16VarP(&mf.Years, "years", "y", 30, "Loan term in years.")

	mortgageCmd.MarkFlagRequired("amount")
	mortgageCmd.MarkFlagRequired("rate")

	// optional flags
	mortgageCmd.Flags().Float64Var(&mf.ExtraMonthlyPayment, "extra-monthly", 0, "Extra monthly payment.")
	mortgageCmd.Flags().Float64Var(&mf.ExtraAnnualPayment, "extra-annual", 0, "Extra annual payment.")
	mortgageCmd.Flags().BoolVar(&mf.MonthlySchedule, "monthly-schedule", false, "Print the monthly amortization schedule.")
	mortgageCmd.Flags().BoolVar(&mf.AnnualSchedule, "annual-schedule", false, "Print the annual amortization schedule.")

	// technically not mutually exclusive, but would need extra logic.
	mortgageCmd.MarkFlagsMutuallyExclusive("extra-monthly", "extra-annual")
	mortgageCmd.MarkFlagsMutuallyExclusive("monthly-schedule", "annual-schedule")
}
