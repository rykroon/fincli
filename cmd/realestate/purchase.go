package realestate

import (
	"github.com/rykroon/fincli/internal/finance"
	"github.com/rykroon/fincli/internal/mortgage"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var purchaseCmd = &cobra.Command{
	Use:   "purchase",
	Short: "Calculate the costs of purchasing a home.",
	Run:   runPurchaseCmd,
}

type purchaseFlags struct {
	Price              finance.Money
	DownPaymentPercent finance.Percent
	Rate               DecimalFlag
	Years              DecimalFlag
	ClosingPercent     finance.Percent
	Escrow             finance.Money
	AnnualTax          finance.Money
	AnnualInsurance    finance.Money
	PmiRate            DecimalFlag
	MonthlyHoa         DecimalFlag
}

func (pf *purchaseFlags) DownPayment() decimal.Decimal {
	return pf.DownPaymentPercent.ApplyTo(pf.Price.Decimal())
}

func (pf *purchaseFlags) LoanAmount() decimal.Decimal {
	return pf.Price.Decimal().Sub(pf.DownPayment())
}

func (pf *purchaseFlags) ClosingCosts() decimal.Decimal {
	return pf.ClosingPercent.ApplyTo(pf.Price.Decimal())
}

var pf purchaseFlags

func runPurchaseCmd(cmd *cobra.Command, args []string) {
	fmt := message.NewPrinter(language.English)
	// Print Summary
	fmt.Printf("Home Price: %v\n", pf.Price)
	fmt.Printf("Down Payment (%v): $%v\n", pf.DownPaymentPercent, pf.DownPayment().StringFixed(0))
	fmt.Printf("Loan Amount: $%v\n", pf.LoanAmount().StringFixed(2))
	fmt.Println("")

	// One-Time costs
	fmt.Printf("--- One-Time costs ---\n")
	fmt.Printf("Closing Costs (%v): %v\n", pf.ClosingPercent, pf.ClosingCosts().StringFixed(2))
	fmt.Printf("Escrow Prepaids: $%v\n", pf.Escrow)
	totalUpfront := decimal.Sum(pf.DownPayment(), pf.ClosingCosts(), pf.Escrow.Decimal())
	fmt.Printf("TOTAL UPFRONT: $%v\n", totalUpfront.StringFixed(2))
	fmt.Println("")

	twelve := decimal.NewFromInt(12)
	oneHundred := decimal.NewFromFloat(100)

	// monthly costs
	fmt.Printf("--- Monthly Costs ---\n")
	p := pf.Price.Decimal().Sub(pf.DownPayment())
	i := pf.Rate.Div(twelve).Div(oneHundred)
	n := pf.Years.Mul(twelve)
	monthlyMortgage := mortgage.CalculateMonthlyPayment(p, i, n)
	monthlyTaxes := pf.AnnualTax.Decimal().Div(twelve)
	monthlyInsurance := pf.AnnualInsurance.Decimal().Div(twelve)
	monthlyPMI := pf.LoanAmount().Mul(pf.PmiRate.Decimal).Div(oneHundred).Div(twelve)

	fmt.Printf("Mortgage Payment: $%v\n", monthlyMortgage.StringFixed(2))
	fmt.Printf("Property Tax: $%v\n", monthlyTaxes.StringFixed(2))
	fmt.Printf("Home Insurance: $%v\n", monthlyInsurance.StringFixed(2))
	if pf.MonthlyHoa.GreaterThan(decimal.Zero) {
		fmt.Printf("HOA: $%v\n", pf.MonthlyHoa.StringFixed(2))
	}

	if monthlyPMI.GreaterThan(decimal.Zero) {
		fmt.Printf("PMI: $%v\n", monthlyPMI.StringFixed(2))
	}

	totalMonthlyCost := decimal.Sum(monthlyMortgage, monthlyTaxes, monthlyInsurance, pf.MonthlyHoa.Decimal, monthlyPMI)

	fmt.Printf("TOTAL MONTHLY: %-12v\n", totalMonthlyCost.StringFixed(2))
}

func init() {
	purchaseCmd.Flags().VarP(&pf.Price, "price", "p", "Home price")
	pf.DownPaymentPercent = finance.NewPercentFromInt(20)
	purchaseCmd.Flags().VarP(&pf.DownPaymentPercent, "down", "d", "Down payment percent (default: 20)")
	purchaseCmd.Flags().VarP(&pf.Rate, "rate", "r", "Mortgage interest rate")
	pf.Years = DecimalFlag{decimal.NewFromInt(30)}
	purchaseCmd.Flags().VarP(&pf.Years, "years", "y", "Mortgage term in years (default: 30)")
	pf.ClosingPercent = finance.NewPercentFromInt(3)
	purchaseCmd.Flags().Var(&pf.ClosingPercent, "closing-percent", "Estimated closing costs (% of price, default: 3)")
	purchaseCmd.Flags().Var(&pf.Escrow, "escrows", "Estimate of prepaid escrow costs")
	purchaseCmd.Flags().VarP(&pf.AnnualTax, "taxes", "t", "Annual property taxes")
	purchaseCmd.Flags().VarP(&pf.AnnualInsurance, "insurance", "i", "Annual homeowners insurance")
	purchaseCmd.Flags().Var(&pf.PmiRate, "pmi", "PMI rate")
	purchaseCmd.Flags().Var(&pf.MonthlyHoa, "hoa", "Monthly HOA fee")

	purchaseCmd.MarkFlagRequired("price")
	purchaseCmd.MarkFlagRequired("rate")

}
