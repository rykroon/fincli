package cmd

import (
	"github.com/rykroon/fincli/internal/flagx"
	"github.com/rykroon/fincli/internal/fmtx"
	"github.com/rykroon/fincli/internal/tax"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
)

var taxCmd = &cobra.Command{
	Use:   "tax",
	Short: "Calculate Income Taxes",
	RunE:  runTaxCmd,
}

type taxFlags struct {
	income       decimal.Decimal
	filingStatus string
	year         uint16
	adjustments  decimal.Decimal
}

var tf taxFlags

func runTaxCmd(cmd *cobra.Command, args []string) error {
	prt := fmtx.NewDecimalPrinter(sep)
	prt.Printf("Gross Income: $%.2v\n", tf.income)
	prt.Println("")

	usTaxSystem, ok := tax.UsFederalRegistry[tf.year]
	if !ok {
		panic("tax system not found")
	}

	taxPayer := tax.NewTaxPayer(
		tf.income,
		tax.FilingStatus(tf.filingStatus),
		tax.Adjustment{Label: "Adjustments", Amount: tf.adjustments},
	)

	usTaxResult := usTaxSystem.CalculateTax(taxPayer)
	effectiveTaxRate := usTaxResult.TaxesDue.Div(tf.income)

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
	ficaTaxSystem, ok := tax.FicaRegistry[tf.year]
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

func init() {
	taxCmd.Flags().VarP(flagx.NewDecVal(&tf.income), "income", "i", "Your gross income")
	taxCmd.Flags().StringVarP(&tf.filingStatus, "filing-status", "f", "single", "Your filing status")
	taxCmd.Flags().Uint16VarP(&tf.year, "year", "y", 2025, "Tax year")
	taxCmd.Flags().Var(flagx.NewDecVal(&tf.adjustments), "adjustments", "adjustments (ex: Reitrement Contributions, Student Loan Interest)")
	taxCmd.MarkFlagRequired("income")
}
