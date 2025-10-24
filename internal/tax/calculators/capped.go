package calculators

import "github.com/shopspring/decimal"

type CappedTax struct {
	Rate  decimal.Decimal
	Upper decimal.Decimal
}

func NewCappedTax(rate, upper float64) CappedTax {
	return CappedTax{decimal.NewFromFloat(rate), decimal.NewFromFloat(upper)}
}

func (t CappedTax) CalculateTax(income decimal.Decimal) decimal.Decimal {
	return decimal.Min(income, t.Upper).Mul(t.Rate)
}
