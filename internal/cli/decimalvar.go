package cli

import (
	"github.com/shopspring/decimal"
)

type DecimalVar struct {
	target *decimal.Decimal
}

func NewDecimalVar(ptr *decimal.Decimal) *DecimalVar {
	return &DecimalVar{target: ptr}
}

func (dv *DecimalVar) Set(s string) error {
	d, err := decimal.NewFromString(s)
	if err != nil {
		return err
	}
	*dv.target = d
	return nil
}

func (dv DecimalVar) String() string {
	if dv.target == nil {
		return ""
	}
	return dv.target.String()
}

func (dv DecimalVar) Type() string { return "decimal" }
