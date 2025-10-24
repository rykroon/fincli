package calculators

import "github.com/shopspring/decimal"

type FlatTax struct {
	Rate decimal.Decimal
}

func NewFlatTax(rate float64) FlatTax {
	return FlatTax{Rate: decimal.NewFromFloat(rate)}
}

func (t FlatTax) CalculateTax(income decimal.Decimal) decimal.Decimal {
	return income.Mul(t.Rate)
}
