package taxes

import "github.com/shopspring/decimal"

type ProgressiveTax struct {
	Brackets []Bracket
}

func NewProgressiveTax(brackets ...Bracket) ProgressiveTax {
	return ProgressiveTax{
		Brackets: brackets,
	}
}

func (t ProgressiveTax) GetMarginalBracket(i decimal.Decimal) Bracket {
	for _, bracket := range t.Brackets {
		if bracket.Lower.LessThan(i) && i.LessThan(bracket.Upper) {
			return bracket
		}
	}
	panic("a bracket could not be found")
}

func (c ProgressiveTax) CalculateTax(income decimal.Decimal) decimal.Decimal {
	tax := decimal.Zero

	for _, b := range c.Brackets {
		if income.LessThanOrEqual(b.Lower) {
			break
		}

		tax = tax.Add(b.CalculateTax(income))

		if income.LessThan(b.Upper) {
			break
		}
	}
	return tax
}

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

func (b Bracket) CalculateTax(income decimal.Decimal) decimal.Decimal {
	if income.LessThan(b.Lower) {
		return decimal.Zero
	}
	upper := decimal.Min(income, b.Upper)
	taxable := upper.Sub(b.Lower)
	return taxable.Mul(b.Rate)
}
