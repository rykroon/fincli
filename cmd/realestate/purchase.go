package realestate

import (
	"github.com/rykroon/fincli/internal/flag"
	"github.com/rykroon/fincli/internal/format"
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
	Price              flag.DecimalFlag
	DownPaymentPercent flag.PercentFlag
	Rate               flag.PercentFlag
	Years              flag.DecimalFlag
	ClosingPercent     flag.PercentFlag
	Escrow             flag.DecimalFlag
	AnnualTax          flag.DecimalFlag
	AnnualInsurance    flag.DecimalFlag
	PmiRate            flag.DecimalFlag
	MonthlyHoa         flag.DecimalFlag
}

func (pf *purchaseFlags) DownPayment() decimal.Decimal {
	return pf.Price.Mul(pf.DownPaymentPercent.Decimal)
}

func (pf *purchaseFlags) LoanAmount() decimal.Decimal {
	return pf.Price.Sub(pf.DownPayment())
}

func (pf *purchaseFlags) ClosingCosts() decimal.Decimal {
	return pf.Price.Mul(pf.ClosingPercent.Decimal)
}

var pf purchaseFlags

func runPurchaseCmd(cmd *cobra.Command, args []string) {
	fmt := message.NewPrinter(language.English)
	// Print Summary
	fmt.Printf("Home Price: %v\n", format.FormatMoney(pf.Price.Decimal))
	fmt.Printf("Down Payment (%v): %v\n", format.FormatPercent(pf.DownPaymentPercent.Decimal), format.FormatMoney(pf.DownPayment()))
	fmt.Printf("Loan Amount: %v\n", format.FormatMoney(pf.LoanAmount()))
	fmt.Println("")

	// One-Time costs
	fmt.Printf("--- One-Time costs ---\n")
	fmt.Printf("Closing Costs (%v): %v\n", format.FormatPercent(pf.ClosingPercent.Decimal), format.FormatMoney(pf.ClosingCosts()))
	fmt.Printf("Escrow Prepaids: %v\n", format.FormatMoney(pf.Escrow.Decimal))
	totalUpfront := decimal.Sum(pf.DownPayment(), pf.ClosingCosts(), pf.Escrow.Decimal)
	fmt.Printf("TOTAL UPFRONT: %v\n", format.FormatMoney(totalUpfront))
	fmt.Println("")

	twelve := decimal.NewFromInt(12)

	// monthly costs
	fmt.Printf("--- Monthly Costs ---\n")
	p := pf.Price.Sub(pf.DownPayment())
	i := pf.Rate.Div(twelve)
	n := pf.Years.Mul(twelve)
	monthlyMortgage := mortgage.CalculateMonthlyPayment(p, i, n)
	monthlyTaxes := pf.AnnualTax.Div(twelve)
	monthlyInsurance := pf.AnnualInsurance.Div(twelve)
	monthlyPMI := pf.LoanAmount().Mul(pf.PmiRate.Decimal).Div(twelve)

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
	pf.DownPaymentPercent = flag.PercentFlag{decimal.NewFromFloat(.2)}
	purchaseCmd.Flags().VarP(&pf.DownPaymentPercent, "down", "d", "Down payment percent (default: 20)")
	purchaseCmd.Flags().VarP(&pf.Rate, "rate", "r", "Mortgage interest rate")
	pf.Years = flag.DecimalFlag{decimal.NewFromInt(30)}
	purchaseCmd.Flags().VarP(&pf.Years, "years", "y", "Mortgage term in years (default: 30)")
	pf.ClosingPercent = flag.PercentFlag{decimal.NewFromFloat(0.03)}
	purchaseCmd.Flags().Var(&pf.ClosingPercent, "closing-percent", "Estimated closing costs (% of price, default: 3)")
	purchaseCmd.Flags().Var(&pf.Escrow, "escrows", "Estimate of prepaid escrow costs")
	purchaseCmd.Flags().VarP(&pf.AnnualTax, "taxes", "t", "Annual property taxes")
	purchaseCmd.Flags().VarP(&pf.AnnualInsurance, "insurance", "i", "Annual homeowners insurance")
	purchaseCmd.Flags().Var(&pf.PmiRate, "pmi", "PMI rate")
	purchaseCmd.Flags().Var(&pf.MonthlyHoa, "hoa", "Monthly HOA fee")

	purchaseCmd.MarkFlagRequired("price")
	purchaseCmd.MarkFlagRequired("rate")

}
