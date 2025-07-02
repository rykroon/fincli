package realestate

import (
	"github.com/rykroon/fincli/internal/mortgage"
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
	Price              float64
	DownPaymentPercent uint8
	Rate               float64
	Years              uint8
	ClosingPercent     uint8
	Escrow             float64
	AnnualTax          float64
	AnnualInsurance    float64
	PmiRate            float64
	MonthlyHoa         float64
}

func (pf *purchaseFlags) DownPayment() float64 {
	return pf.Price * float64(pf.DownPaymentPercent) / 100
}

func (pf *purchaseFlags) LoanAmount() float64 {
	return pf.Price - pf.DownPayment()
}

func (pf *purchaseFlags) ClosingCosts() float64 {
	return pf.Price * float64(pf.ClosingPercent) / 100
}

var pf purchaseFlags

func runPurchaseCmd(cmd *cobra.Command, args []string) {
	fmt := message.NewPrinter(language.English)
	// Print Summary
	fmt.Printf("Home Price: $%.2f\n", pf.Price)
	fmt.Printf("Down Payment (%d%%): $%.2f\n", pf.DownPaymentPercent, pf.DownPayment())
	fmt.Printf("Loan Amount: $%.2f\n", pf.LoanAmount())
	fmt.Println("")

	// One-Time costs
	fmt.Printf("--- One-Time costs ---\n")
	fmt.Printf("Closing Costs (%d%%): %.2f\n", int(pf.ClosingPercent), pf.ClosingCosts())
	fmt.Printf("Escrow Prepaids: $%.2f\n", pf.Escrow)
	totalUpfront := pf.DownPayment() + pf.ClosingCosts() + pf.Escrow
	fmt.Printf("TOTAL UPFRONT: $%.2f\n", totalUpfront)
	fmt.Println("")

	// monthly costs
	fmt.Printf("--- Monthly Costs ---\n")
	p := pf.Price - pf.DownPayment()
	i := pf.Rate / 12 / 100
	n := int(pf.Years) * 12
	monthlyMortgage := mortgage.CalculateMonthlyPayment(p, i, n)
	monthlyTaxes := pf.AnnualTax / 12
	monthlyInsurance := pf.AnnualInsurance / 12
	monthlyPMI := pf.LoanAmount() * pf.PmiRate / 100 / 12

	fmt.Printf("Mortgage Payment: $%.2f\n", monthlyMortgage)
	fmt.Printf("Property Tax: $%.2f\n", monthlyTaxes)
	fmt.Printf("Home Insurance: $%.2f\n", monthlyInsurance)
	if pf.MonthlyHoa > 0 {
		fmt.Printf("HOA: $%.2f\n", pf.MonthlyHoa)
	}

	if monthlyPMI > 0 {
		fmt.Printf("PMI: $%.2f\n", monthlyPMI)
	}

	totalMonthlyCost := monthlyMortgage + monthlyTaxes + monthlyInsurance + pf.MonthlyHoa + monthlyPMI
	fmt.Printf("TOTAL MONTHLY: %-12.2f\n", totalMonthlyCost)
}

func init() {
	purchaseCmd.Flags().Float64VarP(&pf.Price, "price", "p", 0, "Home price")
	purchaseCmd.Flags().Uint8VarP(&pf.DownPaymentPercent, "down", "d", 20, "Down payment percent (default: 20)")
	purchaseCmd.Flags().Float64VarP(&pf.Rate, "rate", "r", 0, "Mortgage interest rate")
	purchaseCmd.Flags().Uint8VarP(&pf.Years, "years", "y", 30, "Mortgage term in years (default: 30)")
	purchaseCmd.Flags().Uint8Var(&pf.ClosingPercent, "closing-percent", 3, "Estimated closing costs (% of price, default: 3)")
	purchaseCmd.Flags().Float64Var(&pf.Escrow, "escrows", 0, "Estimate of prepaid escrow costs")
	purchaseCmd.Flags().Float64VarP(&pf.AnnualTax, "taxes", "t", 0, "Annual property taxes")
	purchaseCmd.Flags().Float64VarP(&pf.AnnualInsurance, "insurance", "i", 0, "Annual homeowners insurance")
	purchaseCmd.Flags().Float64Var(&pf.PmiRate, "pmi", 0, "PMI rate")
	purchaseCmd.Flags().Float64Var(&pf.MonthlyHoa, "hoa", 0, "Monthly HOA fee")

	purchaseCmd.MarkFlagRequired("price")
	purchaseCmd.MarkFlagRequired("rate")

}
