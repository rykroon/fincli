package mortgage

import (
	"github.com/rykroon/fincli/internal/cli"
	"github.com/rykroon/fincli/internal/flag"
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
	monthlyRate := mf.Rate.Div(decimal.NewFromInt(12))
	numPeriods := mf.Years * 12
	sched := mortgage.CalculateSchedule(mf.Amount, monthlyRate, numPeriods, mortgage.NoExtraPayment())

	cmd.Println("Monthly Payment: ", cli.FormatMoney(sched.MonthlyPayment, sep))
	cmd.Println("Total Amount Paid: ", cli.FormatMoney(sched.TotalAmount, sep))
	cmd.Println("Total Interest Paid: ", cli.FormatMoney(sched.TotalInterest, sep))
}

func init() {
	monthlyCmd.Flags().VarP(
		flag.NewDecVal(&mf.Amount), "amount", "a", "The loan amount borrowed.",
	)
	monthlyCmd.Flags().VarP(flag.NewPercentVal(&mf.Rate), "rate", "r", "Annual interest rate.")
	monthlyCmd.Flags().Int64VarP(&mf.Years, "years", "y", 30, "Loan term in years")

	monthlyCmd.MarkFlagRequired("amount")
	monthlyCmd.MarkFlagRequired("rate")

	monthlyCmd.Flags().SortFlags = false
	monthlyCmd.Flags().PrintDefaults()
}
