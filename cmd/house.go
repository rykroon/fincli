package cmd

import (
	"fmt"
	"strings"

	"github.com/rykroon/fincli/internal/flagx"
	"github.com/rykroon/fincli/internal/mortgage"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
)

type houseFlags struct {
	Price              decimal.Decimal
	DownPaymentPercent decimal.Decimal
	Rate               decimal.Decimal
	Years              uint16
	ClosingPercent     decimal.Decimal
	AnnualTax          decimal.Decimal
	AnnualInsurance    decimal.Decimal
	PmiRate            decimal.Decimal
	MonthlyHoa         decimal.Decimal
}

func (hf houseFlags) DownPayment() decimal.Decimal {
	return hf.Price.Mul(hf.DownPaymentPercent)
}

func (hf houseFlags) LoanAmount() decimal.Decimal {
	return hf.Price.Sub(hf.DownPayment())
}

func (hf houseFlags) ClosingCosts() decimal.Decimal {
	return hf.Price.Mul(hf.ClosingPercent)
}

func NewHouseCmd() *cobra.Command {
	var hf houseFlags

	cmd := &cobra.Command{
		Use:   "house",
		Short: "Calculate the costs of purchasing a house.",
		Run: func(cmd *cobra.Command, args []string) {
			runHouseCmd(hf)
		},
	}

	hf.Price = decimal.Zero
	cmd.Flags().VarP(flagx.NewDecimalFlag(&hf.Price), "price", "p", "Home price")

	hf.DownPaymentPercent = decimal.NewFromFloat(.2)
	cmd.Flags().VarP(
		flagx.NewPercentFlag(&hf.DownPaymentPercent),
		"down",
		"d",
		"Down payment percent",
	)

	hf.Rate = decimal.Zero
	cmd.Flags().VarP(
		flagx.NewPercentFlag(&hf.Rate), "rate", "r", "Mortgage interest rate",
	)

	cmd.Flags().Uint16VarP(&hf.Years, "years", "y", 30, "Mortgage term in years")

	hf.ClosingPercent = decimal.NewFromFloat(.03)
	cmd.Flags().Var(
		flagx.NewPercentFlag(&hf.ClosingPercent),
		"closing-percent",
		"Estimated closing costs as a percent",
	)

	hf.AnnualTax = decimal.Zero
	cmd.Flags().VarP(
		flagx.NewDecimalFlag(&hf.AnnualTax), "taxes", "t", "Annual property taxes",
	)

	hf.AnnualInsurance = decimal.Zero
	cmd.Flags().VarP(
		flagx.NewDecimalFlag(&hf.AnnualInsurance),
		"insurance",
		"i",
		"Annual homeowners insurance",
	)

	hf.PmiRate = decimal.Zero
	cmd.Flags().Var(flagx.NewPercentFlag(&hf.PmiRate), "pmi", "PMI rate")

	hf.MonthlyHoa = decimal.Zero
	cmd.Flags().Var(flagx.NewDecimalFlag(&hf.MonthlyHoa), "hoa", "Monthly HOA fee")

	cmd.MarkFlagRequired("price")
	cmd.MarkFlagRequired("rate")

	cmd.Flags().SortFlags = false
	cmd.Flags().PrintDefaults()

	return cmd
}

func runHouseCmd(hf houseFlags) {
	oneHundred := decimal.NewFromInt(100)
	// Print Summary
	prt.Printf("%-20s $%12.2v\n", "Home Price:", hf.Price)
	prt.Printf("%-20s $%12.2v\n", "Loan Amount:", hf.LoanAmount())
	prt.Println("")

	// One-Time costs
	prt.Println("One-Time costs")
	prt.Println(strings.Repeat("-", 20))

	prt.Printf(
		"%-20s $%12.2v\n",
		fmt.Sprintf("Down Payment (%v%%):", hf.DownPaymentPercent.Mul(oneHundred)),
		hf.DownPayment(),
	)
	prt.Printf(
		"%-20s $%12.2v\n",
		fmt.Sprintf("Closing Costs (%v%%):", hf.ClosingPercent.Mul(oneHundred)),
		hf.ClosingCosts(),
	)
	totalUpfront := decimal.Sum(hf.DownPayment(), hf.ClosingCosts())
	prt.Printf("%-20s $%12.2v\n", "Total Upfront:", totalUpfront)
	prt.Println("")

	// monthly costs
	prt.Println("Monthly Costs")
	prt.Println(strings.Repeat("-", 20))
	p := hf.Price.Sub(hf.DownPayment())
	twelve := decimal.NewFromInt(12)
	i := hf.Rate.Div(twelve)
	n := hf.Years * 12
	monthlyMortgage := mortgage.CalculateMonthlyPayment(p, i, n)
	monthlyTaxes := hf.AnnualTax.Div(twelve)
	monthlyInsurance := hf.AnnualInsurance.Div(twelve)
	monthlyPMI := hf.LoanAmount().Mul(hf.PmiRate).Div(twelve)

	prt.Printf("%-20s $%12.2v\n", "Mortgage Payment:", monthlyMortgage)
	prt.Printf("%-20s $%12.2v\n", "Property Tax:", monthlyTaxes)
	prt.Printf("%-20s $%12.2v\n", "Home Insurance:", monthlyInsurance)
	if hf.MonthlyHoa.GreaterThan(decimal.Zero) {
		prt.Printf("%-20s $%12.2v\n", "HOA:", hf.MonthlyHoa)
	}

	if monthlyPMI.GreaterThan(decimal.Zero) {
		prt.Printf("%-20s $%12.2v\n", "PMI:", monthlyPMI)
	}

	totalMonthlyCost := decimal.Sum(
		monthlyMortgage, monthlyTaxes, monthlyInsurance, hf.MonthlyHoa, monthlyPMI,
	)
	prt.Printf("%-20s $%12.2v\n", "Total Monthly:", totalMonthlyCost)
}
