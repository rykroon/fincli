package cmd

import (
	"github.com/rykroon/fincli/internal/cli"
	"github.com/rykroon/fincli/internal/mortgage"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
)

var homePurchCmd = &cobra.Command{
	Use:   "home",
	Short: "Calculate the costs of purchasing a home.",
	Run:   runPurchaseCmd,
}

type purchaseFlags struct {
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

func (pf purchaseFlags) DownPayment() decimal.Decimal {
	return pf.Price.Mul(pf.DownPaymentPercent)
}

func (pf purchaseFlags) LoanAmount() decimal.Decimal {
	return pf.Price.Sub(pf.DownPayment())
}

func (pf purchaseFlags) ClosingCosts() decimal.Decimal {
	return pf.Price.Mul(pf.ClosingPercent)
}

var pf purchaseFlags

func runPurchaseCmd(cmd *cobra.Command, args []string) {
	// Print Summary
	cmd.Println("Home Price: ", cli.FormatMoney(pf.Price, sep))
	cmd.Printf("Down Payment (%s): %s\n", cli.FormatPercent(pf.DownPaymentPercent, 0), cli.FormatMoney(pf.DownPayment(), sep))
	cmd.Println("Loan Amount: ", cli.FormatMoney(pf.LoanAmount(), sep))
	cmd.Println("")

	// One-Time costs
	cmd.Println("--- One-Time costs ---")
	cmd.Printf("Closing Costs (%s): %s\n", cli.FormatPercent(pf.ClosingPercent, 0), cli.FormatMoney(pf.ClosingCosts(), sep))
	totalUpfront := decimal.Sum(pf.DownPayment(), pf.ClosingCosts())
	cmd.Println("Total Upfront: ", cli.FormatMoney(totalUpfront, sep))
	cmd.Println("")

	// monthly costs
	cmd.Println("--- Monthly Costs ---")
	p := pf.Price.Sub(pf.DownPayment())
	twelve := decimal.NewFromInt(12)
	i := pf.Rate.Div(twelve)
	n := pf.Years * 12
	monthlyMortgage := mortgage.CalculateMonthlyPayment(p, i, n)
	monthlyTaxes := pf.AnnualTax.Div(twelve)
	monthlyInsurance := pf.AnnualInsurance.Div(twelve)
	monthlyPMI := pf.LoanAmount().Mul(pf.PmiRate).Div(twelve)

	cmd.Println("Mortgage Payment: ", cli.FormatMoney(monthlyMortgage, sep))
	cmd.Println("Property Tax: ", cli.FormatMoney(monthlyTaxes, sep))
	cmd.Println("Home Insurance: ", cli.FormatMoney(monthlyInsurance, sep))
	if pf.MonthlyHoa.GreaterThan(decimal.Zero) {
		cmd.Println("HOA: ", cli.FormatMoney(pf.MonthlyHoa, sep))
	}

	if monthlyPMI.GreaterThan(decimal.Zero) {
		cmd.Println("PMI: ", cli.FormatMoney(monthlyPMI, sep))
	}

	totalMonthlyCost := decimal.Sum(
		monthlyMortgage, monthlyTaxes, monthlyInsurance, pf.MonthlyHoa, monthlyPMI,
	)
	cmd.Println("Total Monthly: ", cli.FormatMoney(totalMonthlyCost, sep))
}

func init() {
	homePurchCmd.Flags().VarP(cli.DecimalValue(&pf.Price), "price", "p", "Home price")

	pf.DownPaymentPercent = decimal.NewFromFloat(.2)
	homePurchCmd.Flags().VarP(cli.PercentValue(&pf.DownPaymentPercent), "down", "d", "Down payment percent")

	homePurchCmd.Flags().VarP(cli.PercentValue(&pf.Rate), "rate", "r", "Mortgage interest rate")

	homePurchCmd.Flags().Int64VarP(&pf.Years, "years", "y", 30, "Mortgage term in years")

	pf.ClosingPercent = decimal.NewFromFloat(.03)
	homePurchCmd.Flags().Var(cli.PercentValue(&pf.ClosingPercent), "closing-percent", "Estimated closing costs as a percent")
	homePurchCmd.Flags().VarP(cli.DecimalValue(&pf.AnnualTax), "taxes", "t", "Annual property taxes")
	homePurchCmd.Flags().VarP(cli.DecimalValue(&pf.AnnualInsurance), "insurance", "i", "Annual homeowners insurance")
	homePurchCmd.Flags().Var(cli.DecimalValue(&pf.PmiRate), "pmi", "PMI rate")
	homePurchCmd.Flags().Var(cli.DecimalValue(&pf.MonthlyHoa), "hoa", "Monthly HOA fee")

	homePurchCmd.MarkFlagRequired("price")
	homePurchCmd.MarkFlagRequired("rate")

	homePurchCmd.Flags().SortFlags = false
	homePurchCmd.Flags().PrintDefaults()
}
