package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
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

	flagx.DecimalVarP(
		cmd.Flags(),
		&mf.Principal,
		"principal",
		"p",
		decimal.Zero,
		"The principal (loan amount)",
	)

	flagx.PercentVarP(
		cmd.Flags(),
		&mf.Rate,
		"rate",
		"r",
		decimal.Zero,
		"Annual interest rate",
	)

	cmd.Flags().Uint16VarP(&mf.Years, "years", "y", 30, "Loan term in years")

	cmd.MarkFlagRequired("principal")
	cmd.MarkFlagRequired("rate")

	// optional flags
	flagx.DecimalVar(
		cmd.Flags(),
		&mf.ExtraMonthlyPayment,
		"extra-monthly",
		decimal.Zero,
		"Extra monthly payment",
	)

	flagx.DecimalVar(
		cmd.Flags(),
		&mf.ExtraAnnualPayment,
		"extra-annual",
		decimal.Zero,
		"Extra annual payment",
	)

	cmd.Flags().BoolVar(
		&mf.MonthlySchedule,
		"monthly",
		false,
		"Print the monthly amortization schedule",
	)

	cmd.Flags().BoolVar(
		&mf.AnnualSchedule,
		"annual",
		false,
		"Print the annual amortization schedule",
	)

	cmd.MarkFlagsMutuallyExclusive("annual", "monthly")

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
	MonthlySchedule     bool
	AnnualSchedule      bool
}

func (mf mortgageFlags) HasExtraPayment() bool {
	return (mf.ExtraAnnualPayment.GreaterThan(decimal.Zero) ||
		mf.ExtraMonthlyPayment.GreaterThan(decimal.Zero))
}

func runMortgageCmd(mf mortgageFlags) {
	loan := mortgage.NewLoan(mf.Principal, mf.Rate, mf.Years)
	sched := mortgage.CalculateSchedule(loan)
	monthlyPayment := mortgage.CalculateMonthlyPayment(
		loan.Principal, loan.MonthlyRate(), loan.NumPeriods(),
	)

	result := ""

	result += prt.Sprintf("Monthly Payment: $%.2v\n", monthlyPayment)
	if !monthlyPayment.Round(2).Equal(sched.AverageMonthlyPayment().Round(2)) {
		prt.Printf("Average Monthly Payment: $%.2v\n", sched.AverageMonthlyPayment())
	}

	result += prt.Sprintf("Total Amount Paid: $%.2v\n", sched.TotalAmount)
	result += prt.Sprintf("Total Interest Paid: $%.2v\n", sched.TotalInterest)

	twelve := decimal.NewFromInt(12)
	years := sched.NumPeriods().Div(twelve)
	months := sched.NumPeriods().Mod(twelve)
	result += prt.Sprintf("Pay off in %v years and %v months\n", years, months)
	result += prt.Sprintln("")

	if mf.AnnualSchedule {
		result += printAnnualSchedule(sched)
	} else if mf.MonthlySchedule {
		result += printMonthlySchedule(sched)
	} else {
		fmt.Println(result)
		return
	}

	// pipe result into pager

	pager := os.Getenv("PAGER")
	if pager == "" {
		pager = "less" // fallback to less
	}

	cmd := exec.Command(pager)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// We pipe our output *into* the pager's stdin
	pagerIn, _ := cmd.StdinPipe()

	if err := cmd.Start(); err != nil {
		fmt.Println(result) // fallback if pager fails
		return
	}

	io.WriteString(pagerIn, result)
	pagerIn.Close()

	cmd.Wait()
}

func printMonthlySchedule(schedule mortgage.Schedule) string {
	result := ""
	for _, payment := range schedule.Payments {
		if payment.Period%12 == 1 {
			result += prt.Sprintf(
				"%-6s %-12s %-12s %-12s %-12s\n",
				"Month",
				"Principal",
				"Interest",
				"Total",
				"Balance",
			)
			result += prt.Sprintln(strings.Repeat("-", 60))
		}

		result += prt.Sprintf(
			"%-6d $%-11.2v $%-11.2v $%-11.2v $%-11.2v\n",
			payment.Period,
			payment.Principal,
			payment.Interest,
			payment.Total(),
			payment.Balance,
		)

		if payment.Period%12 == 0 {
			result += prt.Sprintf("\t--- End of Year %d ---\n\n", payment.Period/12)
		}
	}
	return result
}

func printAnnualSchedule(schedule mortgage.Schedule) string {
	result := ""
	result += prt.Sprintf(
		"%-6s %-12s %-12s %-12s %-12s\n",
		"Year",
		"Principal",
		"Interest",
		"Total",
		"Balance",
	)
	result += prt.Sprintln(strings.Repeat("-", 60))
	annualPrincipal := decimal.Zero
	annualInterest := decimal.Zero
	annualPayments := decimal.Zero

	for _, payment := range schedule.Payments {
		annualPrincipal = annualPrincipal.Add(payment.Principal)
		annualInterest = annualInterest.Add(payment.Interest)
		annualPayments = annualPayments.Add(payment.Total())

		if payment.Period%12 == 0 {
			result += prt.Sprintf(
				"%-6d $%-11.2v $%-11.2v $%-11.2v $%-11.2v\n",
				payment.Period/12,
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
	return result
}
