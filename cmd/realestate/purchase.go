package realestate

import (
	"github.com/rykroon/fincli/internal/cli"
	"github.com/rykroon/fincli/internal/format"
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
	cmd.Println("Home Price: ", format.FormatMoney(pf.Price))
	cmd.Printf("Down Payment (%v): %v\n", pf.DownPaymentPercent, format.FormatMoney(pf.DownPayment()))
	cmd.Println("Loan Amount: ", format.FormatMoney(pf.LoanAmount()))
	cmd.Println("")

	// One-Time costs
	cmd.Println("--- One-Time costs ---")
	cmd.Printf("Closing Costs (%v): %v\n", pf.ClosingPercent, format.FormatMoney(pf.ClosingCosts()))
	cmd.Println("Escrow Prepaids: ", format.FormatMoney(pf.Escrow))
	totalUpfront := decimal.Sum(pf.DownPayment(), pf.ClosingCosts(), pf.Escrow)
	cmd.Println("TOTAL UPFRONT: ", format.FormatMoney(totalUpfront))
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

	cmd.Println("Mortgage Payment: ", format.FormatMoney(monthlyMortgage))
	cmd.Println("Property Tax: ", format.FormatMoney(monthlyTaxes))
	cmd.Println("Home Insurance: ", format.FormatMoney(monthlyInsurance))
	if pf.MonthlyHoa.GreaterThan(decimal.Zero) {
		cmd.Println("HOA: ", format.FormatMoney(pf.MonthlyHoa))
	}

	if monthlyPMI.GreaterThan(decimal.Zero) {
		cmd.Println("PMI: ", format.FormatMoney(monthlyPMI))
	}

	totalMonthlyCost := decimal.Sum(monthlyMortgage, monthlyTaxes, monthlyInsurance, pf.MonthlyHoa, monthlyPMI)

	cmd.Printf("TOTAL MONTHLY: %-12v\n", format.FormatMoney(totalMonthlyCost))
}

func init() {
	purchaseCmd.Flags().VarP(cli.NewDecimalVar(&pf.Price), "price", "p", "Home price")
	pf.DownPaymentPercent = decimal.NewFromInt(20)
	purchaseCmd.Flags().VarP(cli.NewDecimalVar(&pf.DownPaymentPercent), "down", "d", "Down payment percent (default: 20)")
	purchaseCmd.Flags().VarP(cli.NewDecimalVar(&pf.Rate), "rate", "r", "Mortgage interest rate")
	pf.Years = decimal.NewFromInt(30)
	purchaseCmd.Flags().VarP(cli.NewDecimalVar(&pf.Years), "years", "y", "Mortgage term in years (default: 30)")
	pf.ClosingPercent = decimal.NewFromInt(3)
	purchaseCmd.Flags().Var(cli.NewDecimalVar(&pf.ClosingPercent), "closing-percent", "Estimated closing costs (% of price, default: 3)")
	purchaseCmd.Flags().Var(cli.NewDecimalVar(&pf.Escrow), "escrows", "Estimate of prepaid escrow costs")
	purchaseCmd.Flags().VarP(cli.NewDecimalVar(&pf.AnnualTax), "taxes", "t", "Annual property taxes")
	purchaseCmd.Flags().VarP(cli.NewDecimalVar(&pf.AnnualInsurance), "insurance", "i", "Annual homeowners insurance")
	purchaseCmd.Flags().Var(cli.NewDecimalVar(&pf.PmiRate), "pmi", "PMI rate")
	purchaseCmd.Flags().Var(cli.NewDecimalVar(&pf.MonthlyHoa), "hoa", "Monthly HOA fee")

	purchaseCmd.MarkFlagRequired("price")
	purchaseCmd.MarkFlagRequired("rate")

}
