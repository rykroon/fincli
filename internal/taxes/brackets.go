package taxes

import "github.com/shopspring/decimal"

type Bracket struct {
	Min  decimal.Decimal
	Max  decimal.Decimal
	Rate decimal.Decimal
}

func NewBracket(min, max int64, rate float64) Bracket {
	return Bracket{
		Min:  decimal.NewFromInt(min),
		Max:  decimal.NewFromInt(max),
		Rate: decimal.NewFromFloat(rate),
	}
}

func (b Bracket) CalculateTax(income decimal.Decimal) decimal.Decimal {
	if income.LessThan(b.Min) {
		return decimal.Zero
	}
	upper := b.Max
	if upper.IsZero() || income.LessThan(upper) {
		upper = income
	}

	taxable := upper.Sub(b.Min)
	return taxable.Mul(b.Rate)
}
