package mortgage

import (
	"fmt"
	"strings"

	"github.com/rykroon/fincli/internal/flagx"
	"github.com/rykroon/fincli/internal/fmtx"
	"github.com/rykroon/fincli/internal/mortgage"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
)

func NewAmortizeCmd() *cobra.Command {
	var af amortizeFlags

	cmd := &cobra.Command{
		Use:   "amortize",
		Short: "Print an Amortization Schedule",
		Run: func(cmd *cobra.Command, args []string) {
			sep, _ := flagx.GetRune(cmd.Flags(), "sep")
			prt := fmtx.NewDecimalPrinter(sep)
			runAmortizeCmd(af, prt)
		},
	}

	flagx.DecimalVarP(cmd.Flags(), &af.Principal, "principal", "p", decimal.Zero, "The principal (loan amount)")

	flagx.PercentVarP(cmd.Flags(), &af.Rate, "rate", "r", decimal.Zero, "Annual interest rate")
	cmd.Flags().Int64VarP(&af.Years, "years", "y", 30, "Loan term in years")

	cmd.MarkFlagRequired("principal")
	cmd.MarkFlagRequired("rate")

	// optional flags
	flagx.DecimalVar(cmd.Flags(), &af.ExtraMonthlyPayment, "extra-monthly", decimal.Zero, "Extra monthly payment")
	flagx.DecimalVar(cmd.Flags(), &af.ExtraAnnualPayment, "extra-annual", decimal.Zero, "Extra annual payment")

	cmd.Flags().BoolVar(&af.MonthlySchedule, "monthly-schedule", false, "Print the monthly amortization schedule")
	cmd.Flags().BoolVar(&af.AnnualSchedule, "annual-schedule", false, "Print the annual amortization schedule")

	cmd.MarkFlagsMutuallyExclusive("monthly-schedule", "annual-schedule")

	cmd.Flags().SortFlags = false
	cmd.Flags().PrintDefaults()

	return cmd

}

type amortizeFlags struct {
	Principal           decimal.Decimal
	Rate                decimal.Decimal
	Years               int64
	ExtraMonthlyPayment decimal.Decimal
	ExtraAnnualPayment  decimal.Decimal
	MonthlySchedule     bool
	AnnualSchedule      bool
}

func (af amortizeFlags) HasExtraPayment() bool {
	return af.ExtraAnnualPayment.GreaterThan(decimal.Zero) || af.ExtraMonthlyPayment.GreaterThan(decimal.Zero)
}

func runAmortizeCmd(af amortizeFlags, prt fmtx.DecimalPrinter) {
	loan := mortgage.NewLoan(af.Principal, af.Rate, af.Years)
	sched := mortgage.CalculateSchedule(loan)
	monthlyPayment := mortgage.CalculateMonthlyPayment(loan.Principal, loan.MonthlyRate(), loan.NumPeriods())

	prt.Printf("Monthly Payment: $%.2v\n", monthlyPayment)
	if !monthlyPayment.Round(2).Equal(sched.AverageMonthlyPayment().Round(2)) {
		prt.Printf("Average Monthly Payment: $%.2v\n", sched.AverageMonthlyPayment())
	}

	prt.Printf("Total Amount Paid: $%.2v\n", sched.TotalAmount)
	prt.Printf("Total Interest Paid: $%.2v\n", sched.TotalInterest)

	twelve := decimal.NewFromInt(12)
	years := sched.NumPeriods().Div(twelve)
	months := sched.NumPeriods().Mod(twelve)
	prt.Printf("Pay off in %v years and %v months\n", years, months)
	prt.Println("")

	printMonthlySchedule(sched)

	if af.AnnualSchedule {
		printAnnualSchedule(sched)
	}
}

func printMonthlySchedule(schedule mortgage.Schedule) {
	fmt.Printf(
		"%-6s %-12s %-12s %-12s %-12s\n",
		"Month",
		"Principal",
		"Interest",
		"Total",
		"Balance",
	)
	fmt.Println(strings.Repeat("-", 89))
	for _, payment := range schedule.Payments {
		fmt.Printf(
			"%-6d $%-11s $%-14s $%-14s $%-11s\n",
			payment.Period,
			payment.Principal.StringFixed(2),
			payment.Interest.StringFixed(2),
			payment.Total().StringFixed(2),
			payment.Balance.StringFixed(2),
		)

		if payment.Period%12 == 0 {
			fmt.Printf("\t--- End of Year %d ---\n", payment.Period/12)
		}
	}
}

func printAnnualSchedule(schedule mortgage.Schedule) {
	fmt.Printf(
		"%-6s %-12s %-12s %-12s %-12s %-12s\n",
		"Year",
		"Principal",
		"Total Principal",
		"Interest",
		"Total",
		"Balance",
	)
	fmt.Println(strings.Repeat("-", 89))
	annualPrincipal := decimal.Zero
	annualInterest := decimal.Zero
	annualPayments := decimal.Zero

	for _, payment := range schedule.Payments {
		annualPrincipal = annualPrincipal.Add(payment.Principal)
		annualInterest = annualInterest.Add(payment.Interest)
		annualPayments = annualPayments.Add(payment.Total())

		if payment.Period%12 == 0 {
			fmt.Printf(
				"%-6d $%-11s $%-11s $%-12s $%-11s\n",
				payment.Period/12,
				annualPrincipal.StringFixed(2),
				annualInterest.StringFixed(2),
				annualPayments.StringFixed(2),
				payment.Balance.StringFixed(2),
			)
			annualPrincipal = decimal.Zero
			annualInterest = decimal.Zero
			annualPayments = decimal.Zero
		}
	}
}
