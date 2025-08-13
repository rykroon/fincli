package cli

import (
	"strings"

	"github.com/shopspring/decimal"
)

type decimalValue struct {
	target *decimal.Decimal
}

func DecimalValue(d *decimal.Decimal) *decimalValue {
	return &decimalValue{target: d}
}

func (dv decimalValue) String() string {
	if dv.target == nil {
		return ""
	}
	return dv.target.String()
}

func (dv *decimalValue) Set(s string) error {
	s = strings.ReplaceAll(s, "_", "")
	d, err := decimal.NewFromString(s)
	if err != nil {
		return err
	}
	*dv.target = d
	return nil
}

func (dv decimalValue) Type() string { return "decimal" }
