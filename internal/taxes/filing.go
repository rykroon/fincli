package taxes

import "github.com/shopspring/decimal"

type FilingConfig struct {
	Brackets          []Bracket
	StandardDeduction decimal.Decimal
}

func (c *FilingConfig) setStandardDeductionInt(i int64) {
	c.StandardDeduction = decimal.NewFromInt(i)
}

func (c *FilingConfig) addBracket(b Bracket) {
	c.Brackets = append(c.Brackets, b)
}

func (c FilingConfig) GetMarginalBracket(i decimal.Decimal) Bracket {
	taxableIncome := i.Sub(c.StandardDeduction)
	for _, bracket := range c.Brackets {
		if bracket.Max.IsZero() {
			return bracket
		}
		if bracket.Min.LessThan(taxableIncome) && taxableIncome.LessThan(bracket.Max) {
			return bracket
		}
	}
	panic("a bracket should have been found")
}

func (c FilingConfig) CalculateTax(income decimal.Decimal) decimal.Decimal {
	tax := decimal.Zero
	taxableIncome := income.Sub(c.StandardDeduction)

	for _, b := range c.Brackets {
		if taxableIncome.LessThanOrEqual(b.Min) {
			break
		}

		tax = tax.Add(b.CalculateTax(taxableIncome))

		if taxableIncome.LessThan(b.Max) {
			break
		}
	}
	return tax
}
