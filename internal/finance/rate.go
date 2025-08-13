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
	decimal.Decimal
}

func (r Rate) String() string {
	return r.Decimal.Mul(decimal.NewFromInt(100)).StringFixed(2) + "%"
}

func (r *Rate) Set(s string) error {
	d, err := decimal.NewFromString(s)
	if err != nil {
		return err
	}
	r.Decimal = d.Div(decimal.NewFromInt(100))
	return nil
}

func (r *Rate) Type() string {
	return "rate"
}

func NewRateFromDecimal(d decimal.Decimal) Rate {
	return Rate{d.Div(decimal.NewFromInt(100))}
}

func (r Rate) ApplyTo(d decimal.Decimal) decimal.Decimal {
	return d.Mul(r.Decimal)
}
