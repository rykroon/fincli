package realestate

import (
	"fmt"
	"strings"

	"github.com/rykroon/fincli/internal/flag"
	"github.com/rykroon/fincli/internal/mortgage"
	"github.com/shopspring/decimal"
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
	Amount              flag.DecimalFlag
	Rate                flag.DecimalFlag
	Years               flag.DecimalFlag
	ExtraMonthlyPayment flag.DecimalFlag
	ExtraAnnualPayment  flag.DecimalFlag
	MonthlySchedule     bool
	AnnualSchedule      bool
}

func (mf *mortgageFlags) HasExtraPayment() bool {
	return mf.ExtraAnnualPayment.GreaterThan(decimal.Zero) || mf.ExtraMonthlyPayment.GreaterThan(decimal.Zero)
}

func (mf *mortgageFlags) ExtraPaymentStrategy() mortgage.ExtraPaymentStrategy {
	// return mortgage.PrincipalMatchInterest()
	if mf.ExtraMonthlyPayment.GreaterThan(decimal.Zero) && mf.ExtraAnnualPayment.GreaterThan(decimal.Zero) {
		return mortgage.ExtraMonthlyAndAnnualPayment(mf.ExtraMonthlyPayment.Decimal, mf.ExtraAnnualPayment.Decimal)
	} else if mf.ExtraMonthlyPayment.GreaterThan(decimal.Zero) {
		return mortgage.ExtraMonthlyPayment(mf.ExtraMonthlyPayment.Decimal)
	} else if mf.ExtraAnnualPayment.GreaterThan(decimal.Zero) {
		return mortgage.ExtraAnnualPayment(mf.ExtraAnnualPayment.Decimal)
	} else {
		return mortgage.NoExtraPayment()
	}
}

var mf mortgageFlags

func runMortgageCmd(cmd *cobra.Command, args []string) {
	twelve := decimal.NewFromInt(12)
	i := mf.Rate.Div(twelve).Div(decimal.NewFromInt(100))
	n := mf.Years.Mul(twelve)
	schedule := mortgage.CalculateSchedule(mf.Amount.Decimal, i, n, mf.ExtraPaymentStrategy())

	if !mf.HasExtraPayment() {
		fmt.Printf("Monthly Payment: %v\n", schedule.MonthlyPayment.StringFixed(2))
	}

	fmt.Printf("Total Amount Paid: $%v\n", schedule.TotalAmount.StringFixed(2))
	fmt.Printf("Total Interest Paid: $%v\n", schedule.TotalInterest.StringFixed(2))
	fmt.Printf("Pay off in %v years and %v month(s)\n", schedule.NumPeriods().Div(twelve).StringFixed(0), schedule.NumPeriods().Mod(twelve))
	fmt.Println("")

	if mf.MonthlySchedule {
		printMonthlySchedule(schedule)

	} else if mf.AnnualSchedule {
		printAnnualSchedule(schedule)
	}
}

func printMonthlySchedule(schedule mortgage.Schedule) {
	fmt := message.NewPrinter(language.English)
	fmt.Printf(
		"%-6s %-12s %-12s %-12s %-12s %-12s %-12s\n",
		"Month",
		"Principal",
		"Extra Principal",
		"Total Principal",
		"Interest",
		"Total",
		"Balance",
	)
	fmt.Println(strings.Repeat("-", 89))
	for _, payment := range schedule.Payments {
		fmt.Printf(
			"%-6d $%-11.2f $%-14.2f $%-14.2f $%-11.2f $%-12.2f $%-11.2f\n",
			payment.Period,
			payment.Principal,
			payment.ExtraPrincipal,
			payment.TotalPrincipal(),
			payment.Interest,
			payment.Total(),
			payment.Balance,
		)

		if payment.Period%12 == 0 {
			fmt.Printf("\t--- End of Year %d ---\n", payment.Period/12)
		}
	}
}

func printAnnualSchedule(schedule mortgage.Schedule) {
	fmt := message.NewPrinter(language.English)
	fmt.Printf(
		"%-6s %-12s %-12s %-12s %-12s %-12s %-12s\n",
		"Year",
		"Principal",
		"Extra Principal",
		"Total Principal",
		"Interest",
		"Total",
		"Balance",
	)
	fmt.Println(strings.Repeat("-", 89))
	var annualPrincipal decimal.Decimal = decimal.Zero
	var annualExtraPrincipal decimal.Decimal = decimal.Zero
	var annualTotalPrincipal decimal.Decimal = decimal.Zero
	var annualInterest decimal.Decimal = decimal.Zero
	var annualPayments decimal.Decimal = decimal.Zero

	for _, payment := range schedule.Payments {
		annualPrincipal = annualPrincipal.Add(payment.Principal)
		annualExtraPrincipal = annualExtraPrincipal.Add(payment.ExtraPrincipal)
		annualTotalPrincipal = annualTotalPrincipal.Add(payment.TotalPrincipal())
		annualInterest = annualInterest.Add(payment.Interest)
		annualPayments = annualPayments.Add(payment.Total())

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
			annualPrincipal = decimal.Zero
			annualExtraPrincipal = decimal.Zero
			annualTotalPrincipal = decimal.Zero
			annualInterest = decimal.Zero
			annualPayments = decimal.Zero

		}
	}
}

func init() {
	mortgageCmd.Flags().VarP(&mf.Amount, "amount", "a", "The loan amount borrowed.")
	mortgageCmd.Flags().VarP(&mf.Rate, "rate", "r", "Annual interest rate.")
	mf.Years = flag.DecimalFlag{decimal.NewFromInt(30)}
	mortgageCmd.Flags().VarP(&mf.Years, "years", "y", "Loan term in years.")

	mortgageCmd.MarkFlagRequired("amount")
	mortgageCmd.MarkFlagRequired("rate")

	// optional flags
	mf.ExtraMonthlyPayment = flag.DecimalFlag{decimal.Zero}
	mortgageCmd.Flags().Var(&mf.ExtraMonthlyPayment, "extra-monthly", "Extra monthly payment.")

	mf.ExtraAnnualPayment = flag.DecimalFlag{decimal.Zero}
	mortgageCmd.Flags().Var(&mf.ExtraAnnualPayment, "extra-annual", "Extra annual payment.")

	mortgageCmd.Flags().BoolVar(&mf.MonthlySchedule, "monthly-schedule", false, "Print the monthly amortization schedule.")
	mortgageCmd.Flags().BoolVar(&mf.AnnualSchedule, "annual-schedule", false, "Print the annual amortization schedule.")

	// technically not mutually exclusive, but would need extra logic.
	mortgageCmd.MarkFlagsMutuallyExclusive("extra-monthly", "extra-annual")
	mortgageCmd.MarkFlagsMutuallyExclusive("monthly-schedule", "annual-schedule")
}
