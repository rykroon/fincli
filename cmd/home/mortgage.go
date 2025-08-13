package home

import (
	"strings"

	"github.com/rykroon/fincli/internal/cli"
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
	Amount              decimal.Decimal
	Rate                decimal.Decimal
	Years               decimal.Decimal
	ExtraMonthlyPayment decimal.Decimal
	ExtraAnnualPayment  decimal.Decimal
	MonthlySchedule     bool
	AnnualSchedule      bool
}

func (mf *mortgageFlags) HasExtraPayment() bool {
	return mf.ExtraAnnualPayment.GreaterThan(decimal.Zero) || mf.ExtraMonthlyPayment.GreaterThan(decimal.Zero)
}

func (mf *mortgageFlags) ExtraPaymentStrategy() mortgage.ExtraPaymentStrategy {
	// return mortgage.PrincipalMatchInterest()
	if mf.ExtraMonthlyPayment.GreaterThan(decimal.Zero) && mf.ExtraAnnualPayment.GreaterThan(decimal.Zero) {
		return mortgage.ExtraMonthlyAndAnnualPayment(mf.ExtraMonthlyPayment, mf.ExtraAnnualPayment)
	} else if mf.ExtraMonthlyPayment.GreaterThan(decimal.Zero) {
		return mortgage.ExtraMonthlyPayment(mf.ExtraMonthlyPayment)
	} else if mf.ExtraAnnualPayment.GreaterThan(decimal.Zero) {
		return mortgage.ExtraAnnualPayment(mf.ExtraAnnualPayment)
	} else {
		return mortgage.NoExtraPayment()
	}
}

var mf mortgageFlags

func runMortgageCmd(cmd *cobra.Command, args []string) {
	twelve := decimal.NewFromInt(12)
	oneHundred := decimal.NewFromInt(100)
	monthlyRate := mf.Rate.Div(oneHundred).Div(twelve)
	numPeriods := mf.Years.Mul(twelve)
	schedule := mortgage.CalculateSchedule(mf.Amount, monthlyRate, numPeriods, mf.ExtraPaymentStrategy())

	cmd.Printf("Monthly Payment: $%s\n", schedule.MonthlyPayment.StringFixed(2))
	if !schedule.MonthlyPayment.Round(2).Equal(schedule.AverageMonthlyPayment().Round(2)) {
		cmd.Printf("Average Monthly Payment: $%s\n", schedule.AverageMonthlyPayment().StringFixed(2))
	}

	cmd.Printf("Total Amount Paid: $%s\n", schedule.TotalAmount.StringFixed(2))
	cmd.Printf("Total Interest Paid: $%s\n", schedule.TotalInterest.StringFixed(2))
	cmd.Printf("Pay off in %v years and %v month(s)\n", schedule.NumPeriods().Div(twelve).StringFixed(0), schedule.NumPeriods().Mod(twelve))
	cmd.Println("")

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
			"%-6d $%-11s $%-14s $%-14s $%-11s $%-12s $%-11s\n",
			payment.Period(),
			payment.Principal().StringFixed(2),
			payment.ExtraPrincipal().StringFixed(2),
			payment.TotalPrincipal().StringFixed(2),
			payment.Interest().StringFixed(2),
			payment.Total().StringFixed(2),
			payment.Balance().StringFixed(2),
		)

		if payment.Period()%12 == 0 {
			fmt.Printf("\t--- End of Year %d ---\n", payment.Period()/12)
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
	annualPrincipal := decimal.Zero
	annualExtraPrincipal := decimal.Zero
	annualTotalPrincipal := decimal.Zero
	annualInterest := decimal.Zero
	annualPayments := decimal.Zero

	for _, payment := range schedule.Payments {
		annualPrincipal = annualPrincipal.Add(payment.Principal())
		annualExtraPrincipal = annualExtraPrincipal.Add(payment.ExtraPrincipal())
		annualTotalPrincipal = annualTotalPrincipal.Add(payment.TotalPrincipal())
		annualInterest = annualInterest.Add(payment.Interest())
		annualPayments = annualPayments.Add(payment.Total())

		if payment.Period()%12 == 0 {
			fmt.Printf(
				"%-6d $%-11s $%-14s $%-14s $%-11s $%-12s $%-11s\n",
				payment.Period()/12,
				annualPrincipal.StringFixed(2),
				annualExtraPrincipal.StringFixed(2),
				annualTotalPrincipal.StringFixed(2),
				annualInterest.StringFixed(2),
				annualPayments.StringFixed(2),
				payment.Balance().StringFixed(2),
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
	mortgageCmd.Flags().VarP(
		cli.NewFromDecimal(&mf.Amount), "amount", "a", "The loan amount borrowed.",
	)
	mortgageCmd.Flags().VarP(cli.NewFromDecimal(&mf.Rate), "rate", "r", "Annual interest rate.")
	mf.Years = decimal.NewFromInt(30)
	mortgageCmd.Flags().VarP(cli.NewFromDecimal(&mf.Years), "years", "y", "Loan term in years")

	mortgageCmd.MarkFlagRequired("amount")
	mortgageCmd.MarkFlagRequired("rate")

	// optional flags
	mortgageCmd.Flags().Var(cli.NewFromDecimal(&mf.ExtraMonthlyPayment), "extra-monthly", "Extra monthly payment.")
	mortgageCmd.Flags().Var(cli.NewFromDecimal(&mf.ExtraAnnualPayment), "extra-annual", "Extra annual payment.")

	mortgageCmd.Flags().BoolVar(&mf.MonthlySchedule, "monthly-schedule", false, "Print the monthly amortization schedule.")
	mortgageCmd.Flags().BoolVar(&mf.AnnualSchedule, "annual-schedule", false, "Print the annual amortization schedule.")

	mortgageCmd.MarkFlagsMutuallyExclusive("monthly-schedule", "annual-schedule")

	mortgageCmd.Flags().SortFlags = false
	mortgageCmd.Flags().PrintDefaults()
}
