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
	DownPaymentPct  float64
	Rate            float64
	Years           int
	AnnualTax       float64
	AnnualInsurance float64
	MonthlyHoa      float64
}

var mcf monthlyCostFlags

func runMonthlyCostCmd(cmd *cobra.Command, args []string) {
	fmt := message.NewPrinter(language.English)
	downPayment := mcf.Price * mcf.DownPaymentPct / 100
	p := mcf.Price - downPayment
	i := mcf.Rate / 12 / 100
	n := mcf.Years * 12
	monthlyMortgage := mortgage.CalculateMonthlyPayment(p, i, n)
	monthlyTaxes := mcf.AnnualTax / 12
	monthlyInsurance := mcf.AnnualInsurance / 12

	totalMonthlyCost := monthlyMortgage + monthlyTaxes + monthlyInsurance + mcf.MonthlyHoa
	fmt.Printf("Mortgage Payment: $%.2f\n", monthlyMortgage)
	fmt.Printf("Property Tax: $%.2f\n", monthlyTaxes)
	fmt.Printf("Home Insurance: $%.2f\n", monthlyInsurance)
	fmt.Printf("HOA: $%.2f\n", mcf.MonthlyHoa)
	fmt.Printf("%v\n", strings.Repeat("-", 21))
	fmt.Printf("Total: %-12.2f\n", totalMonthlyCost)
}

func init() {
	monthlyCostCmd.Flags().Float64VarP(&mcf.Price, "price", "p", 0, "Home price.")
	monthlyCostCmd.Flags().Float64VarP(&mcf.DownPaymentPct, "down", "d", 0, "Down payment percent.")
	monthlyCostCmd.Flags().Float64VarP(&mcf.Rate, "rate", "r", 0, "Mortgage interest rate.")
	monthlyCostCmd.Flags().IntVarP(&mcf.Years, "years", "y", 30, "Mortgage term.")
	monthlyCostCmd.Flags().Float64VarP(&mcf.AnnualTax, "tax", "t", 0, "Annual tax amount.")
	monthlyCostCmd.Flags().Float64VarP(&mcf.AnnualInsurance, "insurance", "i", 0, "Annual insurance amount.")
	monthlyCostCmd.Flags().Float64Var(&mcf.MonthlyHoa, "hoa", 0, "Monthly HOA fee.")

	monthlyCostCmd.MarkFlagRequired("price")
	monthlyCostCmd.MarkFlagRequired("down")
	monthlyCostCmd.MarkFlagRequired("rate")
}
