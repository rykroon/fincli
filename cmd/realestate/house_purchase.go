package realestate

import (
	"github.com/rykroon/fincli/internal/mortgage"
	"github.com/spf13/cobra"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var housePurchaseCmd = &cobra.Command{
	Use:   "house-purchase",
	Short: "Calculate the costs of owning a home.",
	Run:   ruHousePurchasertCmd,
}

type housePurchaseFlags struct {
	Price              float64
	DownPaymentPercent uint8
	Rate               float64
	Years              uint16
	ClosingPercent     uint8
	Escrow             float64
	AnnualTax          float64
	AnnualInsurance    float64
	PmiRate            float64
	MonthlyHoa         float64
}

func (mcf *housePurchaseFlags) DownPayment() float64 {
	return hpf.Price * float64(hpf.DownPaymentPercent) / 100
}

func (mcf *housePurchaseFlags) LoanAmount() float64 {
	return mcf.Price - mcf.DownPayment()
}

func (mcf *housePurchaseFlags) ClosingCosts() float64 {
	return hpf.Price * float64(hpf.ClosingPercent) / 100
}

var hpf housePurchaseFlags

func ruHousePurchasertCmd(cmd *cobra.Command, args []string) {
	fmt := message.NewPrinter(language.English)
	// Print Summary
	fmt.Printf("Home Price: $%.2f\n", hpf.Price)
	fmt.Printf("Down Payment (%d%%): $%.2f\n", hpf.DownPaymentPercent, hpf.DownPayment())
	fmt.Printf("Loan Amount: $%.2f\n", hpf.LoanAmount())
	fmt.Println("")

	// One-Time costs
	fmt.Printf("--- One-Time costs ---\n")
	fmt.Printf("Closing Costs (%d%%): %.2f\n", int(hpf.ClosingPercent), hpf.ClosingCosts())
	fmt.Printf("Escrow Prepaids: $%.2f\n", hpf.Escrow)
	totalUpfront := hpf.DownPayment() + hpf.ClosingCosts() + hpf.Escrow
	fmt.Printf("TOTAL UPFRONT: $%.2f\n", totalUpfront)
	fmt.Println("")

	// monthly costs
	fmt.Printf("--- Monthly Costs ---\n")
	p := hpf.Price - hpf.DownPayment()
	i := hpf.Rate / 12 / 100
	n := int(hpf.Years * 12)
	monthlyMortgage := mortgage.CalculateMonthlyPayment(p, i, n)
	monthlyTaxes := hpf.AnnualTax / 12
	monthlyInsurance := hpf.AnnualInsurance / 12
	monthlyPMI := hpf.LoanAmount() * hpf.PmiRate / 100 / 12

	fmt.Printf("Mortgage Payment: $%.2f\n", monthlyMortgage)
	fmt.Printf("Property Tax: $%.2f\n", monthlyTaxes)
	fmt.Printf("Home Insurance: $%.2f\n", monthlyInsurance)
	if hpf.MonthlyHoa > 0 {
		fmt.Printf("HOA: $%.2f\n", hpf.MonthlyHoa)
	}

	if monthlyPMI > 0 {
		fmt.Printf("PMI: $%.2f\n", monthlyPMI)
	}

	totalMonthlyCost := monthlyMortgage + monthlyTaxes + monthlyInsurance + hpf.MonthlyHoa + monthlyPMI
	fmt.Printf("TOTAL MONTHLY: %-12.2f\n", totalMonthlyCost)
}

func init() {
	housePurchaseCmd.Flags().Float64VarP(&hpf.Price, "price", "p", 0, "Home price")
	housePurchaseCmd.Flags().Uint8VarP(&hpf.DownPaymentPercent, "down", "d", 20, "Down payment percent (default: 20)")
	housePurchaseCmd.Flags().Float64VarP(&hpf.Rate, "rate", "r", 0, "Mortgage interest rate")
	housePurchaseCmd.Flags().Uint16VarP(&hpf.Years, "years", "y", 30, "Mortgage term in years (default: 30)")
	housePurchaseCmd.Flags().Uint8Var(&hpf.ClosingPercent, "closing-percent", 3, "Estimated closing costs (% of price, default: 3)")
	housePurchaseCmd.Flags().Float64Var(&hpf.Escrow, "escrows", 0, "Estimate of prepaid escrow costs")
	housePurchaseCmd.Flags().Float64VarP(&hpf.AnnualTax, "taxes", "t", 0, "Annual property taxes")
	housePurchaseCmd.Flags().Float64VarP(&hpf.AnnualInsurance, "insurance", "i", 0, "Annual homeowners insurance")
	housePurchaseCmd.Flags().Float64Var(&hpf.PmiRate, "pmi", 0, "PMI rate")
	housePurchaseCmd.Flags().Float64Var(&hpf.MonthlyHoa, "hoa", 0, "Monthly HOA fee")

	housePurchaseCmd.MarkFlagRequired("price")
	housePurchaseCmd.MarkFlagRequired("rate")

}
