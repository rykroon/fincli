package tax

import (
	"errors"
	"fmt"

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

func runIncomeTaxCmd(cmd *cobra.Command, args []string) error {
	config, ok := taxes.UsFederalTaxTable.GetConfig(itf.year, taxes.FilingStatus(itf.filingStatus))
	if !ok {
		return errors.New("tax table not found")
	}
	taxesDue := config.CalculateTax(itf.salary)

	fmt.Println("Taxes due: ", cli.FormatMoney(taxesDue))
	fmt.Println("Percent of Income: ")
	return nil
}

func init() {
	TaxCmd.Flags().VarP(cli.DecimalValue(&itf.salary), "salary", "s", "Your gross salary")
	TaxCmd.Flags().StringVarP(&itf.filingStatus, "filing-status", "f", "single", "Your filing status")
	TaxCmd.Flags().IntVarP(&itf.year, "year", "y", 2025, "Tax year")
	TaxCmd.MarkFlagRequired("salary")
}
