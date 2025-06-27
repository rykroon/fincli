package cmd

import (
	"strings"

	"github.com/rykroon/ry-cli/internal/mortgage"
	"github.com/spf13/cobra"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var mortgageCmd = &cobra.Command{
	Use:   "mortgage",
	Short: "Calculate mortgage costs.",
	Run:   runMortgageCmd,
}

type mortgageFlags struct {
	Amount        float64
	Rate          float64
	Years         int
	Price         float64
	DownPercent   float64
	ExtraPayment  float64
	PrintSchedule bool
}

func (af *mortgageFlags) Principal() float64 {
	if af.Amount > 0 {
		return af.Amount
	} else {
		return af.Price - af.DownPayment()
	}
}

func (af *mortgageFlags) DownPayment() float64 {
	if af.DownPercent > 0 {
		return af.Price * af.DownPercent / 100
	}
	return 0
}

func (af *mortgageFlags) NumPeriods() int {
	return af.Years * 12
}

func (af *mortgageFlags) MonthlyRate() float64 {
	return af.Rate / 12 / 100
}

var mgf mortgageFlags

func runMortgageCmd(cmd *cobra.Command, args []string) {
	fmt := message.NewPrinter(language.English)
	monthlyPayment, payments := mortgage.CalculatePayments(
		mgf.Principal(), mgf.MonthlyRate(), mgf.NumPeriods(), mgf.ExtraPayment,
	)
	monthlyPayment += mgf.ExtraPayment
	totals := mortgage.SumPayments(payments)
	fmt.Printf("Monthly Payment: %.2f\n", monthlyPayment)
	fmt.Printf("Total Payment Amount: $%.2f\n", totals.TotalPayments)
	fmt.Printf("Total Interest Paid: $%.2f\n", totals.TotalInterest)

	if mgf.PrintSchedule {
		yearNum := 0
		for _, payment := range payments {
			if payment.Period%12 == 1 {
				yearNum += 1
				fmt.Println("")
				fmt.Printf("%-18s~~~ Year %d ~~~\n", "", yearNum)
				fmt.Printf(
					"%-6s %-12s %-12s %-12s\n",
					"Period",
					"Principal",
					"Interest",
					"Balance",
				)
				fmt.Println(strings.Repeat("-", 55))
			}

			// add cumulative interest and equity
			fmt.Printf("%-6d $%-11.2f $%-11.2f $%-11.2f\n",
				payment.Period,
				payment.PrincipalPaid,
				payment.InterestPaid,
				payment.Balance,
			)
		}
	}
}

func init() {
	mortgageCmd.Flags().Float64VarP(&mgf.Amount, "amount", "a", 0, "The loan amount borrowed.")
	mortgageCmd.Flags().Float64VarP(&mgf.Rate, "rate", "r", 0, "Annual interest rate.")
	mortgageCmd.Flags().IntVarP(&mgf.Years, "years", "y", 30, "Loan term in years.")

	// price and down percent can be used instead of amount.
	mortgageCmd.Flags().Float64VarP(&mgf.Price, "price", "p", 0, "Home purchase price.")
	mortgageCmd.Flags().Float64VarP(&mgf.DownPercent, "down-percent", "d", 0, "Down payment as a percent.")

	mortgageCmd.MarkFlagsOneRequired("amount", "price")
	mortgageCmd.MarkFlagsRequiredTogether("price", "down-percent")
	mortgageCmd.MarkFlagsMutuallyExclusive("amount", "price")
	mortgageCmd.MarkFlagsMutuallyExclusive("amount", "down-percent")
	mortgageCmd.MarkFlagRequired("rate")

	// optional flags
	mortgageCmd.Flags().Float64Var(&mgf.ExtraPayment, "extra", 0, "Extra payment per month.")
	mortgageCmd.Flags().BoolVar(&mgf.PrintSchedule, "print-schedule", false, "Print the amortization schedule.")
}
