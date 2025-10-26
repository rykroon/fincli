package tax

import "github.com/shopspring/decimal"

type TaxCalculator interface {
	CalculateTax(decimal.Decimal) decimal.Decimal
}

type FlatTax struct {
	Rate decimal.Decimal
}

func NewFlatTax[T Number](rate T) FlatTax {
	return FlatTax{Rate: numToDecimal(rate)}
}

func (t FlatTax) CalculateTax(income decimal.Decimal) decimal.Decimal {
	return income.Mul(t.Rate)
}

type CappedTax struct {
	Rate  decimal.Decimal
	Upper decimal.Decimal
}

func NewCappedTax[T1 Number, T2 Number](rate T1, upper T2) CappedTax {
	return CappedTax{numToDecimal(rate), numToDecimal(upper)}
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

func (t *ProgressiveTax) AddBracket(b bracket) *ProgressiveTax {
	t.Brackets = append(t.Brackets, b)
	return t
}

type bracket struct {
	Lower decimal.Decimal
	Upper decimal.Decimal
	Rate  decimal.Decimal
}

func newBracket[T1 Number, T2 Number, T3 Number](l T1, u T2, r T3) bracket {
	return bracket{
		Lower: numToDecimal(l),
		Upper: numToDecimal(u),
		Rate:  numToDecimal(r),
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
