package taxes

import "github.com/shopspring/decimal"

type Bracket struct {
	Lower decimal.Decimal
	Upper decimal.Decimal
	Rate  decimal.Decimal
}

func NewBracket(lower, upper int64, rate float64) Bracket {
	return Bracket{
		Lower: decimal.NewFromInt(lower),
		Upper: decimal.NewFromInt(upper),
		Rate:  decimal.NewFromFloat(rate),
	}
}

func (b Bracket) calculateTax(income decimal.Decimal) decimal.Decimal {
	if income.LessThan(b.Lower) {
		return decimal.Zero
	}
	upper := b.Upper
	if income.LessThan(upper) {
		upper = income
	}

	taxable := upper.Sub(b.Lower)
	return taxable.Mul(b.Rate)
}
