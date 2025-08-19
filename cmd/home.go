package cmd

import (
	"github.com/rykroon/fincli/internal/cli"
	"github.com/rykroon/fincli/internal/mortgage"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
)

var homeCmd = &cobra.Command{
	Use:   "home",
	Short: "Calculate the costs of purchasing a home.",
	Run:   runHouseCmd,
}

type homeFlags struct {
	Price              decimal.Decimal
	DownPaymentPercent decimal.Decimal
	Rate               decimal.Decimal
	Years              int64
	ClosingPercent     decimal.Decimal
	AnnualTax          decimal.Decimal
	AnnualInsurance    decimal.Decimal
	PmiRate            decimal.Decimal
	MonthlyHoa         decimal.Decimal
}

func (hf homeFlags) DownPayment() decimal.Decimal {
	return hf.Price.Mul(hf.DownPaymentPercent)
}

func (hf homeFlags) LoanAmount() decimal.Decimal {
	return hf.Price.Sub(hf.DownPayment())
}

func (hf homeFlags) ClosingCosts() decimal.Decimal {
	return hf.Price.Mul(hf.ClosingPercent)
}

var hf homeFlags

func runHouseCmd(cmd *cobra.Command, args []string) {
	// Print Summary
	cmd.Println("Home Price: ", cli.FormatMoney(hf.Price, sep))
	cmd.Printf("Down Payment (%s): %s\n", cli.FormatPercent(hf.DownPaymentPercent, 0), cli.FormatMoney(hf.DownPayment(), sep))
	cmd.Println("Loan Amount: ", cli.FormatMoney(hf.LoanAmount(), sep))
	cmd.Println("")

	// One-Time costs
	cmd.Println("--- One-Time costs ---")
	cmd.Printf("Closing Costs (%s): %s\n", cli.FormatPercent(hf.ClosingPercent, 0), cli.FormatMoney(hf.ClosingCosts(), sep))
	totalUpfront := decimal.Sum(hf.DownPayment(), hf.ClosingCosts())
	cmd.Println("Total Upfront: ", cli.FormatMoney(totalUpfront, sep))
	cmd.Println("")

	// monthly costs
	cmd.Println("--- Monthly Costs ---")
	p := hf.Price.Sub(hf.DownPayment())
	twelve := decimal.NewFromInt(12)
	i := hf.Rate.Div(twelve)
	n := hf.Years * 12
	monthlyMortgage := mortgage.CalculateMonthlyPayment(p, i, n)
	monthlyTaxes := hf.AnnualTax.Div(twelve)
	monthlyInsurance := hf.AnnualInsurance.Div(twelve)
	monthlyPMI := hf.LoanAmount().Mul(hf.PmiRate).Div(twelve)

	cmd.Println("Mortgage Payment: ", cli.FormatMoney(monthlyMortgage, sep))
	cmd.Println("Property Tax: ", cli.FormatMoney(monthlyTaxes, sep))
	cmd.Println("Home Insurance: ", cli.FormatMoney(monthlyInsurance, sep))
	if hf.MonthlyHoa.GreaterThan(decimal.Zero) {
		cmd.Println("HOA: ", cli.FormatMoney(hf.MonthlyHoa, sep))
	}

	if monthlyPMI.GreaterThan(decimal.Zero) {
		cmd.Println("PMI: ", cli.FormatMoney(monthlyPMI, sep))
	}

	totalMonthlyCost := decimal.Sum(
		monthlyMortgage, monthlyTaxes, monthlyInsurance, hf.MonthlyHoa, monthlyPMI,
	)
	cmd.Println("Total Monthly: ", cli.FormatMoney(totalMonthlyCost, sep))
}

func init() {
	homeCmd.Flags().VarP(cli.DecimalValue(&hf.Price), "price", "p", "Home price")

	hf.DownPaymentPercent = decimal.NewFromFloat(.2)
	homeCmd.Flags().VarP(cli.PercentValue(&hf.DownPaymentPercent), "down", "d", "Down payment percent")

	homeCmd.Flags().VarP(cli.PercentValue(&hf.Rate), "rate", "r", "Mortgage interest rate")

	homeCmd.Flags().Int64VarP(&hf.Years, "years", "y", 30, "Mortgage term in years")

	hf.ClosingPercent = decimal.NewFromFloat(.03)
	homeCmd.Flags().Var(cli.PercentValue(&hf.ClosingPercent), "closing-percent", "Estimated closing costs as a percent")
	homeCmd.Flags().VarP(cli.DecimalValue(&hf.AnnualTax), "taxes", "t", "Annual property taxes")
	homeCmd.Flags().VarP(cli.DecimalValue(&hf.AnnualInsurance), "insurance", "i", "Annual homeowners insurance")
	homeCmd.Flags().Var(cli.DecimalValue(&hf.PmiRate), "pmi", "PMI rate")
	homeCmd.Flags().Var(cli.DecimalValue(&hf.MonthlyHoa), "hoa", "Monthly HOA fee")

	homeCmd.MarkFlagRequired("price")
	homeCmd.MarkFlagRequired("rate")

	homeCmd.Flags().SortFlags = false
	homeCmd.Flags().PrintDefaults()
}
