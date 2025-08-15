package home

import (
	"github.com/rykroon/fincli/internal/cli"
	"github.com/rykroon/fincli/internal/mortgage"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
)

var purchaseCmd = &cobra.Command{
	Use:   "purchase",
	Short: "Calculate the costs of purchasing a home.",
	Run:   runPurchaseCmd,
}

type purchaseFlags struct {
	Price              decimal.Decimal
	DownPaymentPercent decimal.Decimal
	Rate               decimal.Decimal
	Years              decimal.Decimal
	ClosingPercent     decimal.Decimal
	Escrow             decimal.Decimal
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
	cmd.Println("Escrow Prepaids: ", cli.FormatMoney(pf.Escrow, sep))
	totalUpfront := decimal.Sum(pf.DownPayment(), pf.ClosingCosts(), pf.Escrow)
	cmd.Println("Total Upfront: ", cli.FormatMoney(totalUpfront, sep))
	cmd.Println("")

	// monthly costs
	cmd.Println("--- Monthly Costs ---")
	p := pf.Price.Sub(pf.DownPayment())
	twelve := decimal.NewFromInt(12)
	i := pf.Rate.Div(twelve)
	n := pf.Years.Mul(twelve)
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
	purchaseCmd.Flags().VarP(cli.DecimalValue(&pf.Price), "price", "p", "Home price")
	pf.DownPaymentPercent = decimal.NewFromFloat(.2)
	purchaseCmd.Flags().VarP(cli.PercentValue(&pf.DownPaymentPercent), "down", "d", "Down payment percent")

	purchaseCmd.Flags().VarP(cli.PercentValue(&pf.Rate), "rate", "r", "Mortgage interest rate")
	pf.Years = decimal.NewFromInt(30)
	purchaseCmd.Flags().VarP(cli.DecimalValue(&pf.Years), "years", "y", "Mortgage term in years")
	pf.ClosingPercent = decimal.NewFromFloat(.03)
	purchaseCmd.Flags().Var(cli.PercentValue(&pf.ClosingPercent), "closing-percent", "Estimated closing costs as a percent")
	purchaseCmd.Flags().Var(cli.DecimalValue(&pf.Escrow), "escrows", "Estimate of prepaid escrow costs")
	purchaseCmd.Flags().VarP(cli.DecimalValue(&pf.AnnualTax), "taxes", "t", "Annual property taxes")
	purchaseCmd.Flags().VarP(cli.DecimalValue(&pf.AnnualInsurance), "insurance", "i", "Annual homeowners insurance")
	purchaseCmd.Flags().Var(cli.DecimalValue(&pf.PmiRate), "pmi", "PMI rate")
	purchaseCmd.Flags().Var(cli.DecimalValue(&pf.MonthlyHoa), "hoa", "Monthly HOA fee")

	purchaseCmd.MarkFlagRequired("price")
	purchaseCmd.MarkFlagRequired("rate")

	purchaseCmd.Flags().SortFlags = false
	purchaseCmd.Flags().PrintDefaults()
}
