package flagx

import (
	"strings"

	"github.com/shopspring/decimal"
	"github.com/spf13/pflag"
)

type decimalValue decimal.Decimal

func newDecimalValue(val decimal.Decimal, p *decimal.Decimal) *decimalValue {
	*p = val
	return (*decimalValue)(p)
}

func (dv *decimalValue) Set(s string) error {
	s = strings.ReplaceAll(s, "_", "")
	d, err := decimal.NewFromString(s)
	if err != nil {
		return err
	}
	*dv = decimalValue(d)
	return nil
}

func (dv decimalValue) Type() string { return "decimal" }

func (dv decimalValue) String() string {
	return decimal.Decimal(dv).String()
}

func DecimalVarP(f *pflag.FlagSet, p *decimal.Decimal, name, shorthand string, value decimal.Decimal, usage string) {
	f.VarP(newDecimalValue(value, p), name, shorthand, usage)
}

func DecimalVar(f *pflag.FlagSet, p *decimal.Decimal, name string, value decimal.Decimal, usage string) {
	DecimalVarP(f, p, name, "", value, usage)
}

func DecimalP(f *pflag.FlagSet, name, shorthand string, value decimal.Decimal, usage string) {
	p := new(decimal.Decimal)
	DecimalVarP(f, p, name, shorthand, value, usage)
}

func Decimal(f *pflag.FlagSet, name string, value decimal.Decimal, usage string) {
	p := new(decimal.Decimal)
	DecimalVarP(f, p, name, "", value, usage)
}
