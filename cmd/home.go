package cmd

import (
	"github.com/rykroon/fincli/internal/flagx"
	"github.com/rykroon/fincli/internal/fmtx"
	"github.com/rykroon/fincli/internal/mortgage"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
)

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

func NewHomeCmd() *cobra.Command {
	var hf homeFlags

	cmd := &cobra.Command{
		Use:   "home",
		Short: "Calculate the costs of purchasing a home.",
		Run: func(cmd *cobra.Command, args []string) {
			sep, _ := flagx.GetRune(cmd.Flags(), "sep")
			prt := fmtx.NewDecimalPrinter(sep)
			runHouseCmd(prt, hf)
		},
	}

	flagx.DecimalVarP(cmd.Flags(), &hf.Price, "price", "p", decimal.Zero, "Home price")

	flagx.PercentVarP(
		cmd.Flags(),
		&hf.DownPaymentPercent,
		"down",
		"d",
		decimal.NewFromFloat(.2),
		"Down payment percent",
	)

	flagx.PercentVarP(
		cmd.Flags(), &hf.Rate, "rate", "r", decimal.Zero, "Mortgage interest rate",
	)

	cmd.Flags().Int64VarP(&hf.Years, "years", "y", 30, "Mortgage term in years")

	flagx.PercentVar(
		cmd.Flags(),
		&hf.ClosingPercent,
		"closing-percent",
		decimal.NewFromFloat(.03),
		"Estimated closing costs as a percent",
	)

	flagx.DecimalVarP(
		cmd.Flags(), &hf.AnnualTax, "taxes", "t", decimal.Zero, "Annual property taxes",
	)

	flagx.DecimalVarP(
		cmd.Flags(),
		&hf.AnnualInsurance,
		"insurance",
		"i",
		decimal.Zero,
		"Annual homeowners insurance",
	)

	flagx.PercentVar(cmd.Flags(), &hf.PmiRate, "pmi", decimal.Zero, "PMI rate")
	flagx.DecimalVar(cmd.Flags(), &hf.MonthlyHoa, "hoa", decimal.Zero, "Monthly HOA fee")

	cmd.MarkFlagRequired("price")
	cmd.MarkFlagRequired("rate")

	cmd.Flags().SortFlags = false
	cmd.Flags().PrintDefaults()

	return cmd
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
