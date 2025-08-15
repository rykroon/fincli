package cli

import (
	"strings"

	"github.com/shopspring/decimal"
)

type percentValue struct {
	target *decimal.Decimal
}

var OneHundred = decimal.NewFromInt(100)

func PercentValue(d *decimal.Decimal) *percentValue {
	return &percentValue{target: d}
}

func (dv percentValue) String() string {
	if dv.target == nil {
		return ""
	}
	return FormatPercent(*dv.target, 0)
}

func (dv *percentValue) Set(s string) error {
	s = strings.ReplaceAll(s, "_", "")
	d, err := decimal.NewFromString(s)
	if err != nil {
		return err
	}
	*dv.target = d.Div(OneHundred)
	return nil
}

func (dv percentValue) Type() string { return "percent" }
