package taxes

import (
	"errors"
	"fmt"

	"github.com/rykroon/fincli/internal/finance"
	"github.com/rykroon/fincli/internal/taxes"
	"github.com/spf13/cobra"
)

var IncomeTaxCmd = &cobra.Command{
	Use:   "income",
	Short: "Calculate Federal Income Taxes",
	RunE:  runIncomeTaxCmd,
}

type incomeTaxFlags struct {
	salary       finance.Money
	filingStatus string
	year         int
}

var itf incomeTaxFlags

func runIncomeTaxCmd(cmd *cobra.Command, args []string) error {
	config, ok := taxes.UsFederalTaxTable.GetConfig(itf.year, taxes.FilingStatus(itf.filingStatus))
	if !ok {
		return errors.New("tax table not found")
	}
	taxesDue := config.CalculateTax(itf.salary.Decimal)

	fmt.Println("taxes due: ", finance.Money{taxesDue})
	return nil
}

func init() {
	IncomeTaxCmd.Flags().VarP(&itf.salary, "salary", "s", "Your gross salary")
	IncomeTaxCmd.Flags().StringVarP(&itf.filingStatus, "filing-status", "f", "single", "Your filing status")
	IncomeTaxCmd.Flags().IntVarP(&itf.year, "year", "y", 2025, "Tax year")
	IncomeTaxCmd.MarkFlagRequired("salary")
}
