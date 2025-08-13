package cli

import (
	"strings"

	"github.com/shopspring/decimal"
)

type DecimalValue struct {
	target *decimal.Decimal
}

func NewFromDecimal(d *decimal.Decimal) *DecimalValue {
	return &DecimalValue{target: d}
}

func (dv DecimalValue) String() string {
	if dv.target == nil {
		return ""
	}
	return dv.target.String()
}

func (dv *DecimalValue) Set(s string) error {
	s = strings.ReplaceAll(s, "_", "")
	d, err := decimal.NewFromString(s)
	if err != nil {
		return err
	}
	*dv.target = d
	return nil
}

func (dv DecimalValue) Type() string { return "decimal" }
