package flagx

import (
	"strings"

	"github.com/shopspring/decimal"
	"github.com/spf13/pflag"
)

type DecimalValue decimal.Decimal

func newDecimalValue(val decimal.Decimal, p *decimal.Decimal) *DecimalValue {
	*p = val
	return (*DecimalValue)(p)
}

func (d *DecimalValue) Set(s string) error {
	s = strings.ReplaceAll(s, "_", "")
	v, err := decimal.NewFromString(s)
	if err != nil {
		return err
	}
	*d = DecimalValue(v)
	return nil
}

func (d DecimalValue) Type() string { return "decimal" }

func (d DecimalValue) String() string {
	return decimal.Decimal(d).String()
}

func DecimalVarP(f *pflag.FlagSet, p *decimal.Decimal, name, shorthand string, value decimal.Decimal, usage string) {
	f.VarP(newDecimalValue(value, p), name, shorthand, usage)
}

func DecimalVar(f *pflag.FlagSet, p *decimal.Decimal, name string, value decimal.Decimal, usage string) {
	DecimalVarP(f, p, name, "", value, usage)
}
