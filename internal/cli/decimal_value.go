package cli

import (
	"strings"

	"github.com/shopspring/decimal"
	"github.com/spf13/pflag"
)

type decimalValue struct {
	target *decimal.Decimal
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

func AddDecimalVar(fs *pflag.FlagSet, p *decimal.Decimal, name string, usage string) {
	fs.Var(&decimalValue{target: p}, name, usage)
}

func AddDecimalVarP(fs *pflag.FlagSet, p *decimal.Decimal, name string, shorthand string, usage string) {
	fs.VarP(&decimalValue{target: p}, name, shorthand, usage)
}
