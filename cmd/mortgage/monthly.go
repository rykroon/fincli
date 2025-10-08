package mortgage

import (
	"github.com/rykroon/fincli/internal/flagx"
	"github.com/rykroon/fincli/internal/fmtx"
	"github.com/rykroon/fincli/internal/mortgage"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
)

var monthlyCmd = &cobra.Command{
	Use:   "monthly",
	Short: "Calculate Monthly Mortgage Payment",
	Run:   runMonthlyCmd,
}

type monthlyFlags struct {
	Principal decimal.Decimal
	Rate      decimal.Decimal
	Years     int64
}

var mf monthlyFlags

func runMonthlyCmd(cmd *cobra.Command, args []string) {
	sep := getSep(cmd)
	prt := fmtx.NewDecimalPrinter(sep)

	loan := mortgage.NewLoan(
		mf.Principal, mf.Rate, mf.Years,
	)
	sched := mortgage.CalculateSchedule(loan)
	monthlyPayment := mortgage.CalculateMonthlyPayment(loan.Principal, loan.MonthlyRate(), loan.NumPeriods())
	prt.Printf("Monthly Payment: $%.2v\n", monthlyPayment)
	prt.Printf("Total Amount Paid: $%.2v\n", sched.TotalAmount)
	prt.Printf("Total Interest Paid: $%.2v\n", sched.TotalInterest)
}

func init() {
	monthlyCmd.Flags().VarP(
		flagx.NewDecVal(&mf.Principal), "principal", "p", "Principal (loan amount)",
	)
	monthlyCmd.Flags().VarP(flagx.NewPercentVal(&mf.Rate), "rate", "r", "Annual interest rate.")
	monthlyCmd.Flags().Int64VarP(&mf.Years, "years", "y", 30, "Loan term in years")

	monthlyCmd.MarkFlagRequired("principal")
	monthlyCmd.MarkFlagRequired("rate")

	monthlyCmd.Flags().SortFlags = false
	monthlyCmd.Flags().PrintDefaults()
}
