package finance

import "github.com/shopspring/decimal"

type RateFrequency string

const (
	Annual   RateFrequency = "annual"
	Monthly  RateFrequency = "monthly"
	Biweekly RateFrequency = "biweekly"
	Weekly   RateFrequency = "weekly"
)

type Rate struct {
	value decimal.Decimal
}

func NewRateFromDecimal(d decimal.Decimal) Rate {
	return Rate{value: d.Div(decimal.NewFromInt(100))}
}

func (r *Rate) Decimal() decimal.Decimal {
	return r.value
}

func (r *Rate) ApplyTo(d decimal.Decimal) decimal.Decimal {
	return d.Mul(r.value)
}

func (r *Rate) String() string {
	return r.value.Mul(decimal.NewFromInt(100)).StringFixed(2) + "%"
}
