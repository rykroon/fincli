package flagx

import (
	"strings"

	"github.com/shopspring/decimal"
	"github.com/spf13/pflag"
)

type percentValue decimal.Decimal

func newpercentValue(val decimal.Decimal, p *decimal.Decimal) *percentValue {
	*p = val
	return (*percentValue)(p)
}

func (pv *percentValue) Set(s string) error {
	s = strings.ReplaceAll(s, "_", "")
	d, err := decimal.NewFromString(s)
	if err != nil {
		return err
	}

	*pv = percentValue(d.Div(decimal.NewFromInt(100)))
	return nil
}

func (pv percentValue) Type() string { return "percent" }

func (pv percentValue) String() string {
	return decimal.Decimal(pv).Mul(decimal.NewFromInt(100)).String()
}

func PercentVarP(f *pflag.FlagSet, p *decimal.Decimal, name, shorthand string, value decimal.Decimal, usage string) {
	f.VarP(newpercentValue(value, p), name, shorthand, usage)
}

func PercentVar(f *pflag.FlagSet, p *decimal.Decimal, name string, value decimal.Decimal, usage string) {
	PercentVarP(f, p, name, "", value, usage)
}

func PercentP(f *pflag.FlagSet, name, shorthand string, value decimal.Decimal, usage string) {
	p := new(decimal.Decimal)
	PercentVarP(f, p, name, shorthand, value, usage)
}

func Percent(f *pflag.FlagSet, name string, value decimal.Decimal, usage string) {
	p := new(decimal.Decimal)
	PercentVarP(f, p, name, "", value, usage)
}
