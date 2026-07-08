package cmd

import (
	"fmt"
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
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := mf.Validate(); err != nil {
				return err
			}
			runMortgageCmd(mf)
			return nil
		},
	}

	cmd.Flags().VarP(
		flagx.NewDecimalFlag(&mf.Principal),
		"principal",
		"p",
		"The principal (loan amount)",
	)

	cmd.Flags().VarP(
		flagx.NewPercentFlag(&mf.Rate),
		"rate",
		"r",
		"Annual interest rate",
	)

	cmd.Flags().Uint16VarP(&mf.Years, "years", "y", 30, "Loan term in years")

	cmd.MarkFlagRequired("principal")
	cmd.MarkFlagRequired("rate")

	// optional flags
	cmd.Flags().Var(
		flagx.NewDecimalFlag(&mf.ExtraMonthlyPayment),
		"extra-monthly",
		"Extra monthly payment",
	)

	cmd.Flags().Var(
		flagx.NewDecimalFlag(&mf.ExtraAnnualPayment),
		"extra-annual",
		"Extra annual payment, applied at the start of each loan year",
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

func (mf mortgageFlags) Validate() error {
	if !mf.Principal.IsPositive() {
		return fmt.Errorf("principal must be greater than zero")
	}
	if mf.Rate.IsNegative() {
		return fmt.Errorf("rate must not be negative")
	}
	if mf.Years == 0 {
		return fmt.Errorf("years must be greater than zero")
	}
	if mf.ExtraMonthlyPayment.IsNegative() {
		return fmt.Errorf("extra monthly payment must not be negative")
	}
	if mf.ExtraAnnualPayment.IsNegative() {
		return fmt.Errorf("extra annual payment must not be negative")
	}
	return nil
}

func runMortgageCmd(mf mortgageFlags) {
	loan := mortgage.NewLoan(mf.Principal, mf.Rate, mf.Years)
	hasExtra := mf.ExtraMonthlyPayment.IsPositive() || mf.ExtraAnnualPayment.IsPositive()
	strategy := mortgage.NewDefaultStrategy()
	if hasExtra {
		strategy = mortgage.NewExtraPaymentStrategy(
			mf.ExtraMonthlyPayment, mf.ExtraAnnualPayment,
		)
	}
	sched := mortgage.CalculateSchedule(loan, strategy)
	monthlyPayment := mortgage.CalculateMonthlyPayment(
		loan.Principal, loan.MonthlyRate(), loan.NumPeriods(),
	)

	prt.Println("Loan")
	prt.Println(strings.Repeat("-", 20))
	prt.Printf("%-20s $%12.2v\n", "Principal:", loan.Principal)
	prt.Printf("%-21s %12.2v%%\n", "Interest Rate:", loan.AnnualRate.Mul(decimal.NewFromInt(100)))
	prt.Printf("%-20s %13s\n", "Term:", fmt.Sprintf("%d years", loan.NumYears))
	prt.Println("")

	prt.Println("Payments")
	prt.Println(strings.Repeat("-", 20))
	prt.Printf("%-20s $%12.2v\n", "Monthly Payment:", monthlyPayment)
	if !monthlyPayment.Round(2).Equal(sched.AverageMonthlyPayment().Round(2)) {
		prt.Printf("%-20s $%12.2v\n", "Avg Monthly Payment:", sched.AverageMonthlyPayment())
	}
	prt.Printf("%-20s $%12.2v\n", "Total Paid:", sched.TotalAmount())
	prt.Printf("%-20s $%12.2v\n", "Total Interest:", sched.TotalInterest)
	prt.Printf("%-20s %13s\n", "Payoff Time:", formatMonths(len(sched.Payments)))

	if hasExtra {
		baseline := mortgage.CalculateSchedule(loan, mortgage.NewDefaultStrategy())
		interestSaved := baseline.TotalInterest.Sub(sched.TotalInterest)
		monthsSaved := len(baseline.Payments) - len(sched.Payments)
		prt.Println("")
		prt.Println("Savings (vs. no extra payments)")
		prt.Println(strings.Repeat("-", 20))
		prt.Printf("%-20s $%12.2v\n", "Interest Saved:", interestSaved)
		prt.Printf("%-20s %13s\n", "Time Saved:", formatMonths(monthsSaved))
	}
	prt.Println("")

	if mf.PrintAnnual {
		printAnnualSchedule(sched)
	} else if mf.PrintMonthly {
		printMonthlySchedule(sched)
	}
}

func formatMonths(months int) string {
	return fmt.Sprintf("%d yrs %d mos", months/12, months%12)
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
