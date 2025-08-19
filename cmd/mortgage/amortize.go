package mortgage

import (
	"fmt"
	"strings"

	"github.com/rykroon/fincli/internal/cli"
	"github.com/rykroon/fincli/internal/mortgage"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
)

var amortizeCmd = &cobra.Command{
	Use:   "amortize",
	Short: "Print an Amortization Schedule",
	Run:   runAmortizeCmd,
}

type amortizeFlags struct {
	Amount              decimal.Decimal
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

func (af amortizeFlags) ExtraPaymentStrategy() mortgage.ExtraPaymentStrategy {
	// return mortgage.PrincipalMatchInterest()
	if af.ExtraMonthlyPayment.GreaterThan(decimal.Zero) && af.ExtraAnnualPayment.GreaterThan(decimal.Zero) {
		return mortgage.ExtraMonthlyAndAnnualPayment(af.ExtraMonthlyPayment, af.ExtraAnnualPayment)
	} else if af.ExtraMonthlyPayment.GreaterThan(decimal.Zero) {
		return mortgage.ExtraMonthlyPayment(af.ExtraMonthlyPayment)
	} else if af.ExtraAnnualPayment.GreaterThan(decimal.Zero) {
		return mortgage.ExtraAnnualPayment(af.ExtraAnnualPayment)
	} else {
		return mortgage.NoExtraPayment()
	}
}

var af amortizeFlags

func runAmortizeCmd(cmd *cobra.Command, args []string) {
	twelve := decimal.NewFromInt(12)
	monthlyRate := af.Rate.Div(twelve)
	numPeriods := af.Years * 12
	sched := mortgage.CalculateSchedule(af.Amount, monthlyRate, numPeriods, af.ExtraPaymentStrategy())

	cmd.Println("Monthly Payment: $", cli.FormatDecimal(sched.MonthlyPayment, sep))
	if !sched.MonthlyPayment.Round(2).Equal(sched.AverageMonthlyPayment().Round(2)) {
		cmd.Println("Average Monthly Payment: $", sched.AverageMonthlyPayment())
	}

	cmd.Println("Total Amount Paid: $", cli.FormatDecimal(sched.TotalAmount, sep))
	cmd.Println("Total Interest Paid: $", cli.FormatDecimal(sched.TotalInterest, sep))

	years := sched.NumPeriods().Div(twelve).StringFixed(0)
	months := sched.NumPeriods().Mod(twelve).StringFixed(0)
	cmd.Println("Pay off in ", years, " years and ", months, "month(s)")
	cmd.Println("")

	printMonthlySchedule(sched)

	if af.AnnualSchedule {
		printAnnualSchedule(sched)
	}
}

func printMonthlySchedule(schedule mortgage.Schedule) {
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
	amortizeCmd.Flags().VarP(
		cli.DecimalValue(&af.Amount), "amount", "a", "The loan amount borrowed.",
	)
	amortizeCmd.Flags().VarP(cli.PercentValue(&af.Rate), "rate", "r", "Annual interest rate.")
	amortizeCmd.Flags().Int64VarP(&af.Years, "years", "y", 30, "Loan term in years")

	amortizeCmd.MarkFlagRequired("amount")
	amortizeCmd.MarkFlagRequired("rate")

	// optional flags
	amortizeCmd.Flags().Var(cli.DecimalValue(&af.ExtraMonthlyPayment), "extra-monthly", "Extra monthly payment.")
	amortizeCmd.Flags().Var(cli.DecimalValue(&af.ExtraAnnualPayment), "extra-annual", "Extra annual payment.")

	amortizeCmd.Flags().BoolVar(&af.MonthlySchedule, "monthly-schedule", false, "Print the monthly amortization schedule.")
	amortizeCmd.Flags().BoolVar(&af.AnnualSchedule, "annual-schedule", false, "Print the annual amortization schedule.")

	amortizeCmd.MarkFlagsMutuallyExclusive("monthly-schedule", "annual-schedule")

	amortizeCmd.Flags().SortFlags = false
	amortizeCmd.Flags().PrintDefaults()
}
