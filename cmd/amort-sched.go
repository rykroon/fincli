package cmd

import (
	"fmt"
	"strings"

	"github.com/rykroon/fincli/internal/cli"
	"github.com/rykroon/fincli/internal/mortgage"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
)

var amortSchedCmd = &cobra.Command{
	Use:   "amort-sched",
	Short: "Print an Amortization Schedule",
	Run:   runAmortSchedCmd,
}

type amortSchedFlags struct {
	Amount              decimal.Decimal
	Rate                decimal.Decimal
	Years               int64
	ExtraMonthlyPayment decimal.Decimal
	ExtraAnnualPayment  decimal.Decimal
	MonthlySchedule     bool
	AnnualSchedule      bool
}

func (asf amortSchedFlags) HasExtraPayment() bool {
	return asf.ExtraAnnualPayment.GreaterThan(decimal.Zero) || asf.ExtraMonthlyPayment.GreaterThan(decimal.Zero)
}

func (asf amortSchedFlags) ExtraPaymentStrategy() mortgage.ExtraPaymentStrategy {
	// return mortgage.PrincipalMatchInterest()
	if asf.ExtraMonthlyPayment.GreaterThan(decimal.Zero) && asf.ExtraAnnualPayment.GreaterThan(decimal.Zero) {
		return mortgage.ExtraMonthlyAndAnnualPayment(asf.ExtraMonthlyPayment, asf.ExtraAnnualPayment)
	} else if asf.ExtraMonthlyPayment.GreaterThan(decimal.Zero) {
		return mortgage.ExtraMonthlyPayment(asf.ExtraMonthlyPayment)
	} else if asf.ExtraAnnualPayment.GreaterThan(decimal.Zero) {
		return mortgage.ExtraAnnualPayment(asf.ExtraAnnualPayment)
	} else {
		return mortgage.NoExtraPayment()
	}
}

var asf amortSchedFlags

func runAmortSchedCmd(cmd *cobra.Command, args []string) {
	twelve := decimal.NewFromInt(12)
	monthlyRate := asf.Rate.Div(twelve)
	numPeriods := asf.Years * 12
	sched := mortgage.CalculateSchedule(asf.Amount, monthlyRate, numPeriods, asf.ExtraPaymentStrategy())

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

	if asf.AnnualSchedule {
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
	amortSchedCmd.Flags().VarP(
		cli.DecimalValue(&asf.Amount), "amount", "a", "The loan amount borrowed.",
	)
	amortSchedCmd.Flags().VarP(cli.PercentValue(&asf.Rate), "rate", "r", "Annual interest rate.")
	amortSchedCmd.Flags().Int64VarP(&asf.Years, "years", "y", 30, "Loan term in years")

	amortSchedCmd.MarkFlagRequired("amount")
	amortSchedCmd.MarkFlagRequired("rate")

	// optional flags
	amortSchedCmd.Flags().Var(cli.DecimalValue(&asf.ExtraMonthlyPayment), "extra-monthly", "Extra monthly payment.")
	amortSchedCmd.Flags().Var(cli.DecimalValue(&asf.ExtraAnnualPayment), "extra-annual", "Extra annual payment.")

	amortSchedCmd.Flags().BoolVar(&asf.MonthlySchedule, "monthly-schedule", false, "Print the monthly amortization schedule.")
	amortSchedCmd.Flags().BoolVar(&asf.AnnualSchedule, "annual-schedule", false, "Print the annual amortization schedule.")

	amortSchedCmd.MarkFlagsMutuallyExclusive("monthly-schedule", "annual-schedule")

	amortSchedCmd.Flags().SortFlags = false
	amortSchedCmd.Flags().PrintDefaults()
}
