package flag

import (
	"fmt"

	"github.com/shopspring/decimal"
)

type DecimalFlag struct {
	decimal.Decimal
}

func (df *DecimalFlag) Set(s string) error {
	d, err := decimal.NewFromString(s)
	if err != nil {
		return fmt.Errorf("invalid decimal: %w", err)
	}
	df.Decimal = d
	return nil
}

func (df DecimalFlag) Type() string { return "decimal" }

type PercentFlag struct {
	decimal.Decimal
}

func (pf *PercentFlag) Set(s string) error {
	d, err := decimal.NewFromString(s)
	if err != nil {
		return fmt.Errorf("invalid decimal: %w", err)
	}
	pf.Decimal = d.Div(decimal.NewFromInt(100))
	return nil
}

func (pf PercentFlag) Type() string { return "percent" }
