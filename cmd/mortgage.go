package cmd

import (
	"github.com/rykroon/fincli/internal/cli"
	"github.com/rykroon/fincli/internal/mortgage"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
)

var mortgageCmd = &cobra.Command{
	Use:   "mortgage",
	Short: "Calculate Mortgage Payment",
	Run:   runMortgageCmd,
}

type mortgageFlags struct {
	Amount decimal.Decimal
	Rate   decimal.Decimal
	Years  int64
}

var mf mortgageFlags

func runMortgageCmd(cmd *cobra.Command, args []string) {
	monthlyRate := mf.Rate.Div(decimal.NewFromInt(12))
	numPeriods := mf.Years * 12
	sched := mortgage.CalculateSchedule(mf.Amount, monthlyRate, numPeriods, mortgage.NoExtraPayment())

	cmd.Println("Monthly Payment: ", cli.FormatMoney(sched.MonthlyPayment, sep))
	cmd.Println("Total Amount Paid: ", cli.FormatMoney(sched.TotalAmount, sep))
	cmd.Println("Total Interest Paid: ", cli.FormatMoney(sched.TotalInterest, sep))
}

func init() {
	mortgageCmd.Flags().VarP(
		cli.DecimalValue(&mf.Amount), "amount", "a", "The loan amount borrowed.",
	)
	mortgageCmd.Flags().VarP(cli.PercentValue(&mf.Rate), "rate", "r", "Annual interest rate.")
	mortgageCmd.Flags().Int64VarP(&mf.Years, "years", "y", 30, "Loan term in years")

	mortgageCmd.MarkFlagRequired("amount")
	mortgageCmd.MarkFlagRequired("rate")

	mortgageCmd.Flags().SortFlags = false
	mortgageCmd.Flags().PrintDefaults()
}
