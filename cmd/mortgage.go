package cmd

import (
	"strings"

	"github.com/rykroon/fincli/internal/flagx"
	"github.com/rykroon/fincli/internal/mortgage"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
)

func NewMortgageCmd() *cobra.Command {
	var mf mortgageFlags

	cmd := &cobra.Command{
		Use:   "mortgage",
		Short: "Calculate a mortgage",
		Run: func(cmd *cobra.Command, args []string) {
			runMortgageCmd(mf)
		},
	}

	mf.Principal = decimal.Zero
	cmd.Flags().VarP(
		flagx.NewDecimalFlag(&mf.Principal),
		"principal",
		"p",
		"The principal (loan amount)",
	)

	mf.Rate = decimal.Zero
	cmd.Flags().VarP(
		flagx.NewDecimalFlag(&mf.Rate),
		"rate",
		"r",
		"Annual interest rate",
	)

	cmd.Flags().Uint16VarP(&mf.Years, "years", "y", 30, "Loan term in years")

	cmd.MarkFlagRequired("principal")
	cmd.MarkFlagRequired("rate")

	// optional flags
	mf.ExtraMonthlyPayment = decimal.Zero
	cmd.Flags().Var(
		flagx.NewDecimalFlag(&mf.ExtraMonthlyPayment),
		"extra-monthly",
		"Extra monthly payment",
	)

	mf.ExtraAnnualPayment = decimal.Zero
	cmd.Flags().Var(
		flagx.NewDecimalFlag(&mf.ExtraAnnualPayment),
		"extra-annual",
		"Extra annual payment",
	)

	cmd.Flags().BoolVar(
		&mf.PrintMonthly,
		"print-monthly",
		false,
		"Print the monthly amortization schedule",
	)

	cmd.Flags().BoolVar(
		&mf.PrintAnnual,
		"print-annual",
		false,
		"Print the annual amortization schedule",
	)

	cmd.MarkFlagsMutuallyExclusive("print-annual", "print-monthly")

	cmd.Flags().SortFlags = false
	cmd.Flags().PrintDefaults()

	return cmd

}

type mortgageFlags struct {
	Principal           decimal.Decimal
	Rate                decimal.Decimal
	Years               uint16
	ExtraMonthlyPayment decimal.Decimal
	ExtraAnnualPayment  decimal.Decimal
	PrintMonthly        bool
	PrintAnnual         bool
}

func (mf mortgageFlags) HasExtraPayment() bool {
	return (mf.ExtraAnnualPayment.GreaterThan(decimal.Zero) ||
		mf.ExtraMonthlyPayment.GreaterThan(decimal.Zero))
}

func runMortgageCmd(mf mortgageFlags) {
	loan := mortgage.NewLoan(mf.Principal, mf.Rate, mf.Years)
	strategy := mortgage.NewDefaultStrategy()
	if !mf.ExtraMonthlyPayment.IsZero() {
		strategy = mortgage.NewExtraMonthlyStrategy(mf.ExtraMonthlyPayment)
	} else if !mf.ExtraAnnualPayment.IsZero() {
		strategy = mortgage.NewExtraAnnualStrategy(mf.ExtraAnnualPayment)
	}
	sched := mortgage.CalculateSchedule(loan, strategy)
	monthlyPayment := mortgage.CalculateMonthlyPayment(
		loan.Principal, loan.MonthlyRate(), loan.NumPeriods(),
	)

	prt.Printf("Monthly Payment: $%.2v\n", monthlyPayment)
	if !monthlyPayment.Round(2).Equal(sched.AverageMonthlyPayment().Round(2)) {
		prt.Printf("Average Monthly Payment: $%.2v\n", sched.AverageMonthlyPayment())
	}

	prt.Printf("Total Amount Paid: $%.2v\n", sched.TotalAmount())
	prt.Printf("Total Interest Paid: $%.2v\n", sched.TotalInterest)

	twelve := decimal.NewFromInt(12)
	years := sched.NumPeriods().Div(twelve)
	months := sched.NumPeriods().Mod(twelve)
	prt.Printf("Pay off in %.0v years and %0v months\n", years, months)
	prt.Println("")

	if mf.PrintAnnual {
		printAnnualSchedule(sched)
	} else if mf.PrintMonthly {
		printMonthlySchedule(sched)
	}
}

func printMonthlySchedule(schedule *mortgage.Schedule) {
	for _, payment := range schedule.Payments {
		if payment.Period%12 == 1 {
			prt.Printf(
				"%-6s %-12s %-12s %-12s %-12s\n",
				"Month",
				"Principal",
				"Interest",
				"Total",
				"Balance",
			)
			prt.Println(strings.Repeat("-", 60))
		}

		prt.Printf(
			"%-6d $%-11.2v $%-11.2v $%-11.2v $%-11.2v\n",
			payment.Period,
			payment.Principal,
			payment.Interest,
			payment.Total(),
			payment.Balance,
		)

		if payment.Period%12 == 0 {
			prt.Printf("\t--- End of Year %d ---\n\n", payment.Period/12)
		}
	}
}

func printAnnualSchedule(schedule *mortgage.Schedule) {
	prt.Printf(
		"%-6s %-12s %-12s %-12s %-12s\n",
		"Year",
		"Principal",
		"Interest",
		"Total",
		"Balance",
	)
	prt.Println(strings.Repeat("-", 60))
	annualPrincipal := decimal.Zero
	annualInterest := decimal.Zero
	annualPayments := decimal.Zero

	for idx, payment := range schedule.Payments {
		annualPrincipal = annualPrincipal.Add(payment.Principal)
		annualInterest = annualInterest.Add(payment.Interest)
		annualPayments = annualPayments.Add(payment.Total())

		if payment.Period%12 == 0 || idx == len(schedule.Payments)-1 {
			prt.Printf(
				"%-6d $%-11.2v $%-11.2v $%-11.2v $%-11.2v\n",
				idx/12+1,
				annualPrincipal,
				annualInterest,
				annualPayments,
				payment.Balance,
			)
			annualPrincipal = decimal.Zero
			annualInterest = decimal.Zero
			annualPayments = decimal.Zero
		}
	}
}
