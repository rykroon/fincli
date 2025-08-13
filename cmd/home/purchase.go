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

func (pf *purchaseFlags) DownPayment() decimal.Decimal {
	return pf.Price.Mul(pf.DownPaymentPercent.Div(decimal.NewFromInt(100)))
}

func (pf *purchaseFlags) LoanAmount() decimal.Decimal {
	return pf.Price.Sub(pf.DownPayment())
}

func (pf *purchaseFlags) ClosingCosts() decimal.Decimal {
	return pf.Price.Mul(pf.ClosingPercent.Div(decimal.NewFromInt(100)))
}

var pf purchaseFlags

func runPurchaseCmd(cmd *cobra.Command, args []string) {
	// Print Summary
	cmd.Println("Home Price: ", cli.FormatMoney(pf.Price))
	cmd.Printf("Down Payment (%v): %v\n", pf.DownPaymentPercent, cli.FormatMoney(pf.DownPayment()))
	cmd.Println("Loan Amount: ", cli.FormatMoney(pf.LoanAmount()))
	cmd.Println("")

	// One-Time costs
	cmd.Println("--- One-Time costs ---")
	cmd.Printf("Closing Costs (%v): %v\n", pf.ClosingPercent, cli.FormatMoney(pf.ClosingCosts()))
	cmd.Println("Escrow Prepaids: ", cli.FormatMoney(pf.Escrow))
	totalUpfront := decimal.Sum(pf.DownPayment(), pf.ClosingCosts(), pf.Escrow)
	cmd.Println("TOTAL UPFRONT: ", cli.FormatMoney(totalUpfront))
	cmd.Println("")

	twelve := decimal.NewFromInt(12)

	// monthly costs
	cmd.Println("--- Monthly Costs ---")
	p := pf.Price.Sub(pf.DownPayment())
	i := pf.Rate.Div(twelve).Div(decimal.NewFromInt(100))
	n := pf.Years.Mul(twelve)
	monthlyMortgage := mortgage.CalculateMonthlyPayment(p, i, n)
	monthlyTaxes := pf.AnnualTax.Div(twelve)
	monthlyInsurance := pf.AnnualInsurance.Div(twelve)
	monthlyPMI := pf.LoanAmount().Mul(pf.PmiRate).Div(twelve)

	cmd.Println("Mortgage Payment: ", cli.FormatMoney(monthlyMortgage))
	cmd.Println("Property Tax: ", cli.FormatMoney(monthlyTaxes))
	cmd.Println("Home Insurance: ", cli.FormatMoney(monthlyInsurance))
	if pf.MonthlyHoa.GreaterThan(decimal.Zero) {
		cmd.Println("HOA: ", cli.FormatMoney(pf.MonthlyHoa))
	}

	if monthlyPMI.GreaterThan(decimal.Zero) {
		cmd.Println("PMI: ", cli.FormatMoney(monthlyPMI))
	}

	totalMonthlyCost := decimal.Sum(monthlyMortgage, monthlyTaxes, monthlyInsurance, pf.MonthlyHoa, monthlyPMI)

	cmd.Printf("TOTAL MONTHLY: %-12v\n", cli.FormatMoney(totalMonthlyCost))
}

func init() {
	purchaseCmd.Flags().VarP(cli.DecimalValue(&pf.Price), "price", "p", "Home price")
	pf.DownPaymentPercent = decimal.NewFromInt(20)
	purchaseCmd.Flags().VarP(cli.DecimalValue(&pf.DownPaymentPercent), "down", "d", "Down payment percent")

	purchaseCmd.Flags().VarP(cli.DecimalValue(&pf.Rate), "rate", "r", "Mortgage interest rate")
	pf.Years = decimal.NewFromInt(30)
	purchaseCmd.Flags().VarP(cli.DecimalValue(&pf.Years), "years", "y", "Mortgage term in years")
	pf.ClosingPercent = decimal.NewFromInt(3)
	purchaseCmd.Flags().Var(cli.DecimalValue(&pf.ClosingPercent), "closing-percent", "Estimated closing costs as a percent")
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
