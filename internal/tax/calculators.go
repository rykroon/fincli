package tax

import "github.com/shopspring/decimal"

type TaxCalculator interface {
	CalculateTax(decimal.Decimal) decimal.Decimal
}

type FlatTax struct {
	Rate decimal.Decimal
}

func NewFlatTax(rate float64) FlatTax {
	return FlatTax{Rate: decimal.NewFromFloat(rate)}
}

func (t FlatTax) CalculateTax(income decimal.Decimal) decimal.Decimal {
	return income.Mul(t.Rate)
}

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

type ProgressiveTax struct {
	Brackets []bracket
}

func NewProgressiveTax(brackets ...bracket) ProgressiveTax {
	return ProgressiveTax{
		Brackets: brackets,
	}
}

func (t ProgressiveTax) GetMarginalBracket(income decimal.Decimal) bracket {
	for _, bracket := range t.Brackets {
		if bracket.Lower.LessThan(income) && income.LessThan(bracket.Upper) {
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

		tax = tax.Add(b.calculateTax(income))

		if income.LessThan(b.Upper) {
			break
		}
	}
	return tax
}

func (t *ProgressiveTax) AddBracket(lower, upper int64, rate float64) *ProgressiveTax {
	t.Brackets = append(t.Brackets, newBracket(lower, upper, rate))
	return t
}

type bracket struct {
	Lower decimal.Decimal
	Upper decimal.Decimal
	Rate  decimal.Decimal
}

func newBracket(lower, upper int64, rate float64) bracket {
	return bracket{
		Lower: decimal.NewFromInt(lower),
		Upper: decimal.NewFromInt(upper),
		Rate:  decimal.NewFromFloat(rate),
	}
}

func (b bracket) calculateTax(income decimal.Decimal) decimal.Decimal {
	if income.LessThan(b.Lower) {
		return decimal.Zero
	}
	upper := decimal.Min(income, b.Upper)
	taxable := upper.Sub(b.Lower)
	return taxable.Mul(b.Rate)
}
