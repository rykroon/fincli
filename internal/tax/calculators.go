package tax

import "github.com/shopspring/decimal"

type TaxCalculator interface {
	CalculateTax(decimal.Decimal) decimal.Decimal
}

type FlatTax struct {
	Rate decimal.Decimal `json:"rate"`
}

func (t FlatTax) CalculateTax(income decimal.Decimal) decimal.Decimal {
	return income.Mul(t.Rate)
}

type CappedTax struct {
	Upper decimal.Decimal `json:"upper"`
	Rate  decimal.Decimal `json:"rate"`
}

func (t CappedTax) CalculateTax(income decimal.Decimal) decimal.Decimal {
	return decimal.Min(income, t.Upper).Mul(t.Rate)
}

type ProgressiveTax struct {
	Brackets []Bracket
}

func (t ProgressiveTax) GetMarginalBracket(income decimal.Decimal) Bracket {
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

		tax = tax.Add(b.CalculateTax(income))

		if income.LessThan(b.Upper) {
			break
		}
	}
	return tax
}

type Bracket struct {
	Lower decimal.Decimal `json:"lower"`
	Upper decimal.Decimal `json:"upper"`
	Rate  decimal.Decimal `json:"rate"`
}

func (b Bracket) CalculateTax(income decimal.Decimal) decimal.Decimal {
	if income.LessThan(b.Lower) {
		return decimal.Zero
	}
	upper := decimal.Min(income, b.Upper)
	taxable := upper.Sub(b.Lower)
	return taxable.Mul(b.Rate)
}
