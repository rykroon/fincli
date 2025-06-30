package realestate

import (
	"strings"

	"github.com/rykroon/fincli/internal/mortgage"
	"github.com/spf13/cobra"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var monthlyCostCmd = &cobra.Command{
	Use:   "monthly-costs",
	Short: "Calculate the monthly costs of owning a home.",
	Run:   runMonthlyCostCmd,
}

type monthlyCostFlags struct {
	Price           float64
	DownPaymentAmt  float64
	DownPaymentPct  float64
	Rate            float64
	Years           int
	AnnualTax       float64
	AnnualInsurance float64
	AnnualPMI       float64
	MonthlyHoa      float64
}

func (mcf *monthlyCostFlags) DownPayment() float64 {
	if mcf.DownPaymentPct > 0 {
		return mcf.Price * mcf.DownPaymentPct / 100
	}
	return mcf.DownPaymentAmt
}

var mcf monthlyCostFlags

func runMonthlyCostCmd(cmd *cobra.Command, args []string) {
	fmt := message.NewPrinter(language.English)
	p := mcf.Price - mcf.DownPayment()
	i := mcf.Rate / 12 / 100
	n := mcf.Years * 12
	monthlyMortgage := mortgage.CalculateMonthlyPayment(p, i, n)
	monthlyTaxes := mcf.AnnualTax / 12
	monthlyInsurance := mcf.AnnualInsurance / 12
	monthlyPMI := mcf.AnnualPMI / 12

	fmt.Printf("Mortgage Payment: $%.2f\n", monthlyMortgage)
	fmt.Printf("Property Tax: $%.2f\n", monthlyTaxes)
	fmt.Printf("Home Insurance: $%.2f\n", monthlyInsurance)
	if mcf.MonthlyHoa > 0 {
		fmt.Printf("HOA: $%.2f\n", mcf.MonthlyHoa)
	}

	if mcf.AnnualPMI > 0 {
		fmt.Printf("PMI: $%.2f\n", monthlyPMI)
	}

	fmt.Printf("%v\n", strings.Repeat("-", 21))
	totalMonthlyCost := monthlyMortgage + monthlyTaxes + monthlyInsurance + mcf.MonthlyHoa + monthlyPMI
	fmt.Printf("Total: %-12.2f\n", totalMonthlyCost)
}

func init() {
	monthlyCostCmd.Flags().Float64VarP(&mcf.Price, "price", "p", 0, "Home price.")
	monthlyCostCmd.Flags().Float64VarP(&mcf.DownPaymentAmt, "down-payment", "d", 0, "Down payment amount.")
	monthlyCostCmd.Flags().Float64VarP(&mcf.DownPaymentPct, "down-percent", "D", 0, "Down payment percent.")
	monthlyCostCmd.Flags().Float64VarP(&mcf.Rate, "rate", "r", 0, "Mortgage interest rate.")
	monthlyCostCmd.Flags().IntVarP(&mcf.Years, "years", "y", 30, "Mortgage term.")
	monthlyCostCmd.Flags().Float64VarP(&mcf.AnnualTax, "tax", "t", 0, "Annual tax amount.")
	monthlyCostCmd.Flags().Float64VarP(&mcf.AnnualInsurance, "insurance", "i", 0, "Annual insurance amount.")
	monthlyCostCmd.Flags().Float64Var(&mcf.AnnualPMI, "pmi", 0, "Private mortgage insurance (PMI).")
	monthlyCostCmd.Flags().Float64Var(&mcf.MonthlyHoa, "hoa", 0, "Monthly HOA fee.")

	monthlyCostCmd.MarkFlagRequired("price")
	monthlyCostCmd.MarkFlagRequired("down")
	monthlyCostCmd.MarkFlagRequired("rate")

	monthlyCostCmd.MarkFlagsMutuallyExclusive("down-payment", "down-percent")

}
