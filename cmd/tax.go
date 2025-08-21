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
}

var tf taxFlags

func runTaxCmd(cmd *cobra.Command, args []string) error {
	config, ok := taxes.UsFederalTaxTable.GetConfig(tf.year, taxes.FilingStatus(tf.filingStatus))
	if !ok {
		return errors.New("tax table not found")
	}

	taxesDue := config.CalculateTax(tf.income)
	effectiveTaxRate := taxesDue.Div(tf.income)
	bracket := config.GetBracketByIncome(tf.income)

	prt := fmtx.NewDecimalPrinter(sep)
	prt.Printf("%-18s $%.2v\n", "Taxes Due:", taxesDue)
	prt.Printf("%-18s %.2v%%\n", "Effective Tax Rate:", effectiveTaxRate.Mul(decimal.NewFromInt(100)))
	prt.Printf("Tax Bracket: %12v%%\n", bracket.Rate.Mul(decimal.NewFromInt(100)))
	prt.Printf("Standard Deduction: $%.2v\n", config.StandardDeduction)
	return nil
}

func init() {
	taxCmd.Flags().VarP(flagx.NewDecVal(&tf.income), "income", "i", "Your gross income")
	taxCmd.Flags().StringVarP(&tf.filingStatus, "filing-status", "f", "single", "Your filing status")
	taxCmd.Flags().IntVarP(&tf.year, "year", "y", 2025, "Tax year")
	taxCmd.MarkFlagRequired("income")
}
