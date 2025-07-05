package realestate

import (
	"fmt"

	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
)

var RealEstateCmd = &cobra.Command{
	Use:   "real-estate",
	Short: "Real estate calculation.",
}

func init() {
	RealEstateCmd.AddCommand(purchaseCmd)
	RealEstateCmd.AddCommand(mortgageCmd)
}

type DecimalFlag struct {
	decimal.Decimal
}

func (df *DecimalFlag) Set(s string) error {
	d, err := decimal.NewFromString(s)
	if err != nil {
		return fmt.Errorf("not a valid decimal: %w", err)
	}
	df.Decimal = d
	return nil
}

func (df *DecimalFlag) Type() string {
	return "decimal"
}
