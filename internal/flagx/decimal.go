package flagx

import (
	"strings"

	"github.com/shopspring/decimal"
)

func NewDecimalFlag(p *decimal.Decimal) *Flag[decimal.Decimal] {
	return &Flag[decimal.Decimal]{
		ptr:        p,
		type_:      "decimal",
		stringFunc: func(v decimal.Decimal) string { return v.String() },
		setFunc: func(s string, p *decimal.Decimal) error {
			s = strings.ReplaceAll(s, "_", "")
			v, err := decimal.NewFromString(s)
			if err != nil {
				return err
			}
			*p = v
			return nil
		},
	}
}
