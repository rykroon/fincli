package cmd

import (
	"github.com/rykroon/fincli/internal/flagx"
	"github.com/rykroon/fincli/internal/fmtx"
	"github.com/rykroon/fincli/internal/tax"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
)

func NewTaxCmd() *cobra.Command {
	var income decimal.Decimal
	var filingStatus string
	var year uint16
	var adjustments decimal.Decimal

	runE := func(cmd *cobra.Command, args []string) error {
		sep, _ := flagx.GetRune(cmd.Flags(), "sep")
		prt := fmtx.NewDecimalPrinter(sep)
		taxPayer := tax.NewTaxPayer(
			income,
			tax.FilingStatus(filingStatus),
			tax.Adjustment{Label: "Adjustments", Amount: adjustments},
		)
		runTaxCmd(prt, year, taxPayer)
		return nil
	}

	cmd := &cobra.Command{
		Use:   "tax",
		Short: "Calculate Income Taxes",
		RunE:  runE,
	}
	flagx.DecimalVarP(cmd.Flags(), &income, "income", "i", decimal.Zero, "Your gross income")
	cmd.Flags().StringVarP(&filingStatus, "filing-status", "f", "single", "Your filing status")
	cmd.Flags().Uint16VarP(&year, "year", "y", 2025, "Tax year")
	flagx.DecimalVar(cmd.Flags(), &adjustments, "adjustments", decimal.Zero, "adjustments (ex: Retirement Contributions, Student Loan Interest)")
	cmd.MarkFlagRequired("income")
	return cmd
}

func runTaxCmd(prt fmtx.DecimalPrinter, year uint16, taxPayer tax.TaxPayer) error {
	prt.Printf("Gross Income: $%.2v\n", taxPayer.Income)
	prt.Println("")

	usTaxSystem, ok := tax.UsFederalRegistry[year]
	if !ok {
		panic("tax system not found")
	}

	usTaxResult := usTaxSystem.CalculateTax(taxPayer)
	effectiveTaxRate := usTaxResult.TaxesDue.Div(taxPayer.Income)

	oneHundred := decimal.NewFromInt(100)

	prt.Println("Federal Tax")
	prt.Printf("  Adjusted Gross Income: $%.2v\n", usTaxResult.AdjustedGrossIncome)
	prt.Printf("  Standard Deduction: $%.2v\n", usTaxResult.StandardDeduction)
	prt.Printf("  Taxable Income: $%.2v\n", usTaxResult.TaxableIncome)
	prt.Printf("  Taxes Due: $%.2v\n", usTaxResult.TaxesDue)
	prt.Printf("  Marginal Tax Rate: %v%%\n", usTaxResult.MarginalTaxRate.Mul(oneHundred))
	prt.Printf("  Effective Tax Rate: %.2v%%\n", effectiveTaxRate.Mul(oneHundred))
	prt.Println("")

	// FICA Tax
	ficaTaxSystem, ok := tax.FicaRegistry[year]
	if !ok {
		panic("FICA tax system not found")
	}

	ficaTaxResult := ficaTaxSystem.CalculateTax(taxPayer)

	prt.Println("FICA Tax")
	prt.Printf(
		"  Social Security Tax (%v%%): $%.2v\n",
		ficaTaxSystem.SocialSecurityTax.Rate.Mul(oneHundred),
		ficaTaxResult.SocialSecurityTaxDue,
	)
	prt.Printf(
		"  Medicare Tax (%v%%): $%.2v\n",
		ficaTaxSystem.MedicareTax.Rate.Mul(oneHundred),
		ficaTaxResult.MedicareTaxDue,
	)

	return nil
}
