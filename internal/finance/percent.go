package finance

import "github.com/shopspring/decimal"

type Percent struct {
	value decimal.Decimal
}

func NewPercentFromDecimal(d decimal.Decimal) Percent {
	return Percent{value: d.Div(decimal.NewFromInt(100))}
}

func (p *Percent) Decimal() decimal.Decimal {
	return p.value
}

func (p *Percent) ApplyTo(d decimal.Decimal) decimal.Decimal {
	return d.Mul(p.value)
}

func (p *Percent) String() string {
	return p.value.Mul(decimal.NewFromInt(100)).StringFixed(0) + "%"
}
