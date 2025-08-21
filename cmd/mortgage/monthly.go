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
	Amount decimal.Decimal
	Rate   decimal.Decimal
	Years  int64
}

var mf monthlyFlags

func runMonthlyCmd(cmd *cobra.Command, args []string) {
	sep := getSep(cmd)
	prt := fmtx.NewDecimalPrinter(sep)

	monthlyRate := mf.Rate.Div(decimal.NewFromInt(12))
	numPeriods := mf.Years * 12
	sched := mortgage.CalculateSchedule(mf.Amount, monthlyRate, numPeriods, mortgage.NoExtraPayment())

	prt.Printf("Monthly Payment: $%.2v\n", sched.MonthlyPayment)
	prt.Printf("Total Amount Paid: $%.2v\n", sched.TotalAmount)
	prt.Printf("Total Interest Paid: $%.2v\n", sched.TotalInterest)
}

func init() {
	monthlyCmd.Flags().VarP(
		flagx.NewDecVal(&mf.Amount), "amount", "a", "The loan amount borrowed.",
	)
	monthlyCmd.Flags().VarP(flagx.NewPercentVal(&mf.Rate), "rate", "r", "Annual interest rate.")
	monthlyCmd.Flags().Int64VarP(&mf.Years, "years", "y", 30, "Loan term in years")

	monthlyCmd.MarkFlagRequired("amount")
	monthlyCmd.MarkFlagRequired("rate")

	monthlyCmd.Flags().SortFlags = false
	monthlyCmd.Flags().PrintDefaults()
}
