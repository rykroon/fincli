package finance

import (
	"fmt"

	"github.com/shopspring/decimal"
)

type Percent struct {
	value decimal.Decimal
}

func NewPercentFromDecimal(d decimal.Decimal) Percent {
	return Percent{d.Div(decimal.NewFromInt(100))}
}

func NewPercentFromInt(i int64) Percent {
	return NewPercentFromDecimal(decimal.NewFromInt(i))
}

func (p Percent) ApplyTo(v decimal.Decimal) decimal.Decimal {
	return p.value.Mul(v)
}

func (p Percent) String() string {
	return p.value.Mul(decimal.NewFromInt(100)).StringFixed(0) + "%"
}

func (p *Percent) Set(s string) error {
	d, err := decimal.NewFromString(s)
	if err != nil {
		return fmt.Errorf("invalid decimal: %w", err)
	}
	p.value = d.Div(decimal.NewFromInt(100))
	return nil
}

func (p Percent) Type() string {
	return "percent"
}
