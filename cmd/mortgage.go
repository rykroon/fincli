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
	Years  decimal.Decimal
}

var mf mortgageFlags

func runMortgageCmd(cmd *cobra.Command, args []string) {
	twelve := decimal.NewFromInt(12)
	oneHundred := decimal.NewFromInt(100)
	monthlyRate := mf.Rate.Div(oneHundred).Div(twelve)
	numPeriods := mf.Years.Mul(twelve)
	sched := mortgage.CalculateSchedule(mf.Amount, monthlyRate, numPeriods, mortgage.NoExtraPayment())

	cmd.Println("Monthly Payment: ", cli.FormatMoney(sched.MonthlyPayment, sep))
	cmd.Println("Total Amount Paid: ", cli.FormatMoney(sched.TotalAmount, sep))
	cmd.Println("Total Interest Paid: ", cli.FormatMoney(sched.TotalInterest, sep))
}

func init() {
	mortgageCmd.Flags().VarP(
		cli.DecimalValue(&mf.Amount), "amount", "a", "The loan amount borrowed.",
	)
	mortgageCmd.Flags().VarP(cli.DecimalValue(&mf.Rate), "rate", "r", "Annual interest rate.")
	mf.Years = decimal.NewFromInt(30)
	mortgageCmd.Flags().VarP(cli.DecimalValue(&mf.Years), "years", "y", "Loan term in years")

	mortgageCmd.MarkFlagRequired("amount")
	mortgageCmd.MarkFlagRequired("rate")

	mortgageCmd.Flags().SortFlags = false
	mortgageCmd.Flags().PrintDefaults()
}
