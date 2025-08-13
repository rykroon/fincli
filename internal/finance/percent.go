package finance

import "github.com/shopspring/decimal"

type Percent decimal.Decimal

func NewPercentFromDecimal(d decimal.Decimal) Percent {
	return Percent(d.Div(decimal.NewFromInt(100)))
}

func (p Percent) Decimal() decimal.Decimal {
	return decimal.Decimal(p)
}

func (p Percent) ApplyTo(d decimal.Decimal) decimal.Decimal {
	return d.Mul(p.Decimal())
}

func (p Percent) String() string {
	return p.Decimal().Mul(decimal.NewFromInt(100)).StringFixed(0) + "%"
}
