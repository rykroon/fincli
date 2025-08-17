package cmd

import (
	"errors"

	"github.com/rykroon/fincli/internal/cli"
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

var itf taxFlags

func runTaxCmd(cmd *cobra.Command, args []string) error {
	config, ok := taxes.UsFederalTaxTable.GetConfig(itf.year, taxes.FilingStatus(itf.filingStatus))
	if !ok {
		return errors.New("tax table not found")
	}

	taxesDue := config.CalculateTax(itf.income)
	effectiveTaxRate := taxesDue.Div(itf.income)
	bracket := config.GetBracketByIncome(itf.income)

	cmd.Printf("Taxes Due: \t\t%s\n", cli.FormatMoney(taxesDue, sep))
	cmd.Printf("Effective Tax Rate: \t%s\n", cli.FormatPercent(effectiveTaxRate, 2))
	cmd.Printf("Tax Bracket: \t\t%s\n", cli.FormatPercent(bracket.Rate, 0))
	cmd.Printf("Standard Deduction: \t%s\n", cli.FormatMoney(config.StandardDeduction, sep))
	return nil
}

func init() {
	taxCmd.Flags().VarP(cli.DecimalValue(&itf.income), "income", "i", "Your gross income")
	taxCmd.Flags().StringVarP(&itf.filingStatus, "filing-status", "f", "single", "Your filing status")
	taxCmd.Flags().IntVarP(&itf.year, "year", "y", 2025, "Tax year")
	taxCmd.MarkFlagRequired("income")
}
