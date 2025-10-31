package cmd

import (
	"github.com/rykroon/fincli/internal/flagx"
	"github.com/rykroon/fincli/internal/fmtx"
	"github.com/rykroon/fincli/internal/mortgage"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
)

func NewHomeCmd() *cobra.Command {
	var hf homeFlags

	cmd := &cobra.Command{
		Use:   "home",
		Short: "Calculate the costs of purchasing a home.",
		Run: func(cmd *cobra.Command, args []string) {
			sep := getSep(cmd)
			prt := fmtx.NewDecimalPrinter(sep)
			runHouseCmd(prt, hf)
		},
	}

	cmd.Flags().VarP(flagx.NewDecVal(&hf.Price), "price", "p", "Home price")

	hf.DownPaymentPercent = decimal.NewFromFloat(.2)
	cmd.Flags().VarP(flagx.NewPercentVal(&hf.DownPaymentPercent), "down", "d", "Down payment percent")

	cmd.Flags().VarP(flagx.NewPercentVal(&hf.Rate), "rate", "r", "Mortgage interest rate")

	cmd.Flags().Int64VarP(&hf.Years, "years", "y", 30, "Mortgage term in years")

	hf.ClosingPercent = decimal.NewFromFloat(.03)
	cmd.Flags().Var(flagx.NewPercentVal(&hf.ClosingPercent), "closing-percent", "Estimated closing costs as a percent")
	cmd.Flags().VarP(flagx.NewDecVal(&hf.AnnualTax), "taxes", "t", "Annual property taxes")
	cmd.Flags().VarP(flagx.NewDecVal(&hf.AnnualInsurance), "insurance", "i", "Annual homeowners insurance")
	cmd.Flags().Var(flagx.NewPercentVal(&hf.PmiRate), "pmi", "PMI rate")
	cmd.Flags().Var(flagx.NewDecVal(&hf.MonthlyHoa), "hoa", "Monthly HOA fee")

	cmd.MarkFlagRequired("price")
	cmd.MarkFlagRequired("rate")

	cmd.Flags().SortFlags = false
	cmd.Flags().PrintDefaults()

	return cmd
}

type homeFlags struct {
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

func (hf homeFlags) DownPayment() decimal.Decimal {
	return hf.Price.Mul(hf.DownPaymentPercent)
}

func (hf homeFlags) LoanAmount() decimal.Decimal {
	return hf.Price.Sub(hf.DownPayment())
}

func (hf homeFlags) ClosingCosts() decimal.Decimal {
	return hf.Price.Mul(hf.ClosingPercent)
}

func runHouseCmd(prt fmtx.DecimalPrinter, hf homeFlags) {
	oneHundred := decimal.NewFromInt(100)
	// Print Summary
	prt.Printf("Home Price: $%.2v\n", hf.Price)
	prt.Printf("Down Payment (%v%%): $%.2v\n", hf.DownPaymentPercent.Mul(oneHundred), hf.DownPayment())
	prt.Printf("Loan Amount: $%.2v\n", hf.LoanAmount())
	prt.Println("")

	// One-Time costs
	prt.Println("--- One-Time costs ---")
	prt.Printf("Closing Costs (%v%%): $%.2v\n", hf.ClosingPercent.Mul(oneHundred), hf.ClosingCosts())
	totalUpfront := decimal.Sum(hf.DownPayment(), hf.ClosingCosts())
	prt.Printf("Total Upfront: $%.2v\n", totalUpfront)
	prt.Println("")

	// monthly costs
	prt.Println("--- Monthly Costs ---")
	p := hf.Price.Sub(hf.DownPayment())
	twelve := decimal.NewFromInt(12)
	i := hf.Rate.Div(twelve)
	n := hf.Years * 12
	monthlyMortgage := mortgage.CalculateMonthlyPayment(p, i, n)
	monthlyTaxes := hf.AnnualTax.Div(twelve)
	monthlyInsurance := hf.AnnualInsurance.Div(twelve)
	monthlyPMI := hf.LoanAmount().Mul(hf.PmiRate).Div(twelve)

	prt.Printf("Mortgage Payment: $%.2v\n", monthlyMortgage)
	prt.Printf("Property Tax: $%.2v\n", monthlyTaxes)
	prt.Printf("Home Insurance: $%.2v\n", monthlyInsurance)
	if hf.MonthlyHoa.GreaterThan(decimal.Zero) {
		prt.Printf("HOA: $%.2v\n", hf.MonthlyHoa)
	}

	if monthlyPMI.GreaterThan(decimal.Zero) {
		prt.Printf("PMI: $%.2v\n", monthlyPMI)
	}

	totalMonthlyCost := decimal.Sum(
		monthlyMortgage, monthlyTaxes, monthlyInsurance, hf.MonthlyHoa, monthlyPMI,
	)
	prt.Printf("Total Monthly: $%.2v\n", totalMonthlyCost)
}
