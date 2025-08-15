package tax

import (
	"errors"

	"github.com/rykroon/fincli/internal/cli"
	"github.com/rykroon/fincli/internal/taxes"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
)

var TaxCmd = &cobra.Command{
	Use:   "tax",
	Short: "Calculate Federal Income Taxes",
	RunE:  runIncomeTaxCmd,
}

type taxFlags struct {
	salary       decimal.Decimal
	filingStatus string
	year         int
}

var itf taxFlags

var comma rune = ','

// var underScore rune = '_'
var sep = &comma

func runIncomeTaxCmd(cmd *cobra.Command, args []string) error {
	config, ok := taxes.UsFederalTaxTable.GetConfig(itf.year, taxes.FilingStatus(itf.filingStatus))
	if !ok {
		return errors.New("tax table not found")
	}

	taxesDue := config.CalculateTax(itf.salary)
	effectiveTaxRate := taxesDue.Div(itf.salary)

	cmd.Println("Taxes Due: ", cli.FormatMoney(taxesDue, sep))
	cmd.Println("Effective Tax Rate: ", cli.FormatPercent(effectiveTaxRate, 2))
	return nil
}

func init() {
	TaxCmd.Flags().VarP(cli.DecimalValue(&itf.salary), "salary", "s", "Your gross salary")
	TaxCmd.Flags().StringVarP(&itf.filingStatus, "filing-status", "f", "single", "Your filing status")
	TaxCmd.Flags().IntVarP(&itf.year, "year", "y", 2025, "Tax year")
	TaxCmd.MarkFlagRequired("salary")
}
