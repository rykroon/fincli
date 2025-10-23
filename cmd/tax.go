package cmd

import (
	"errors"

	"github.com/rykroon/fincli/internal/flagx"
	"github.com/rykroon/fincli/internal/fmtx"
	"github.com/rykroon/fincli/internal/taxes"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
)

var taxCmd = &cobra.Command{
	Use:   "tax",
	Short: "Calculate Federal Income Taxes",
	RunE:  runTaxCmd,
}

type taxFlags struct {
	income       decimal.Decimal
	filingStatus string
	year         int
	adjustments  decimal.Decimal
}

var tf taxFlags

func runTaxCmd(cmd *cobra.Command, args []string) error {
	prt := fmtx.NewDecimalPrinter(sep)
	prt.Printf("Gross Income: $%.2v\n", tf.income)
	prt.Println("")

	config, ok := taxes.UsFederalTaxTable.GetConfig(tf.year, taxes.FilingStatus(tf.filingStatus))
	if !ok {
		return errors.New("tax table not found")
	}

	adjustedGrossIncome := tf.income.Sub(tf.adjustments)
	taxesDue := config.CalculateTax(adjustedGrossIncome)
	effectiveTaxRate := taxesDue.Div(tf.income)
	bracket := config.GetMarginalBracket(adjustedGrossIncome)

	oneHundred := decimal.NewFromInt(100)

	prt.Println("Federal Tax")
	prt.Printf("  Adjusted Gross Income: $%.2v\n", adjustedGrossIncome)
	prt.Printf("  Standard Deduction: $%.2v\n", config.StandardDeduction)
	prt.Printf("  Taxable Income: $%.2v\n", adjustedGrossIncome.Sub(config.StandardDeduction))
	prt.Printf("  Taxes Due: $%.2v\n", taxesDue)
	prt.Printf("  Marginal Tax Rate: %v%%\n", bracket.Rate.Mul(oneHundred))
	prt.Printf("  Effective Tax Rate: %.2v%%\n", effectiveTaxRate.Mul(oneHundred))
	prt.Println("")

	// FICA Tax
	socialSecurityTax := taxes.SocialSecurityTax.CalculateTax(tf.income)
	medicareTax := taxes.MedicareTax.CalculateTax(tf.income)
	prt.Println("FICA Tax")
	prt.Printf("  Social Security Tax (%v%%): $%.2v\n", taxes.SocialSecurityTax.Rate.Mul(oneHundred), socialSecurityTax)
	prt.Printf("  Medicare Tax (%v%%): $%.2v\n", taxes.MedicareTax.Rate.Mul(oneHundred), medicareTax)

	return nil
}

func init() {
	taxCmd.Flags().VarP(flagx.NewDecVal(&tf.income), "income", "i", "Your gross income")
	taxCmd.Flags().StringVarP(&tf.filingStatus, "filing-status", "f", "single", "Your filing status")
	taxCmd.Flags().IntVarP(&tf.year, "year", "y", 2025, "Tax year")
	taxCmd.Flags().Var(flagx.NewDecVal(&tf.adjustments), "adjustments", "adjustments (ex: Reitrement Contributions, Student Loan Interest)")
	taxCmd.MarkFlagRequired("income")
}
