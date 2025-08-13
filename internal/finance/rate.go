package finance

import "github.com/shopspring/decimal"

type RateFrequency string

const (
	Annual   RateFrequency = "annual"
	Monthly  RateFrequency = "monthly"
	Biweekly RateFrequency = "biweekly"
	Weekly   RateFrequency = "weekly"
)

type Rate decimal.Decimal

func NewRateFromDecimal(d decimal.Decimal) Rate {
	return Rate(d.Div(decimal.NewFromInt(100)))
}

func (r Rate) Decimal() decimal.Decimal {
	return decimal.Decimal(r)
}

func (r Rate) ApplyTo(d decimal.Decimal) decimal.Decimal {
	return d.Mul(r.Decimal())
}

func (r Rate) String() string {
	return r.Decimal().Mul(decimal.NewFromInt(100)).StringFixed(2) + "%"
}
