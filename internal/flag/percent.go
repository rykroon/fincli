package flag

import (
	"strings"

	"github.com/shopspring/decimal"
)

type PercentVal struct {
	ptr *decimal.Decimal
}

func NewPercentVal(d *decimal.Decimal) *PercentVal {
	return &PercentVal{ptr: d}
}

func (pv PercentVal) String() string {
	if pv.ptr == nil {
		return ""
	}
	return pv.ptr.Mul(decimal.NewFromInt(100)).String()
}

func (pv *PercentVal) Set(s string) error {
	s = strings.ReplaceAll(s, "_", "")
	d, err := decimal.NewFromString(s)
	if err != nil {
		return err
	}
	*pv.ptr = d.Div(decimal.NewFromInt(100))
	return nil
}

func (pv PercentVal) Type() string { return "percent" }
