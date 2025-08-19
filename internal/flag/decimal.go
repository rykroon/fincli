package flag

import (
	"strings"

	"github.com/shopspring/decimal"
)

type DecVal struct {
	ptr *decimal.Decimal
}

func NewDecVal(d *decimal.Decimal) *DecVal {
	return &DecVal{ptr: d}
}

func (dv DecVal) String() string {
	if dv.ptr == nil {
		return ""
	}
	return dv.ptr.String()
}

func (dv *DecVal) Set(s string) error {
	s = strings.ReplaceAll(s, "_", "")
	d, err := decimal.NewFromString(s)
	if err != nil {
		return err
	}
	*dv.ptr = d
	return nil
}

func (dv DecVal) Type() string { return "decimal" }
