package mortgage

import (
	"github.com/rykroon/fincli/internal/flagx"
	"github.com/rykroon/fincli/internal/fmtx"
	"github.com/rykroon/fincli/internal/mortgage"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
)

func NewMontlyCmd() *cobra.Command {
	var principal decimal.Decimal
	var rate decimal.Decimal
	var years int64

	cmd := &cobra.Command{
		Use:   "monthly",
		Short: "Calculate Monthly Mortgage Payment",
		Run: func(cmd *cobra.Command, args []string) {
			sep, _ := flagx.GetRune(cmd.Flags(), "sep")
			prt := fmtx.NewNumberPrinter(sep)
			loan := mortgage.NewLoan(principal, rate, years)
			runMonthlyCmd(loan, prt)
		},
	}

	flagx.DecimalVarP(cmd.Flags(), &principal, "principal", "p", decimal.Zero, "Principal (loan amount)")
	flagx.PercentVarP(cmd.Flags(), &rate, "rate", "r", decimal.Zero, "Annual interest rate")
	cmd.Flags().Int64VarP(&years, "years", "y", 30, "Loan term in years")

	cmd.MarkFlagRequired("principal")
	cmd.MarkFlagRequired("rate")

	cmd.Flags().SortFlags = false
	cmd.Flags().PrintDefaults()
	return cmd
}

func runMonthlyCmd(loan *mortgage.Loan, prt fmtx.NumberPrinter) {
	sched := mortgage.CalculateSchedule(loan)
	monthlyPayment := mortgage.CalculateMonthlyPayment(loan.Principal, loan.MonthlyRate(), loan.NumPeriods())
	prt.Printf("Monthly Payment: $%.2v\n", monthlyPayment)
	prt.Printf("Total Amount Paid: $%.2v\n", sched.TotalAmount)
	prt.Printf("Total Interest Paid: $%.2v\n", sched.TotalInterest)
}
