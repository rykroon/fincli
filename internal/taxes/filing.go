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
		if bracket.Lower.LessThan(taxableIncome) && taxableIncome.LessThan(bracket.Upper) {
			return bracket
		}
	}
	panic("a bracket could not be found")
}

func (c FilingConfig) CalculateTax(income decimal.Decimal) decimal.Decimal {
	tax := decimal.Zero
	taxableIncome := income.Sub(c.StandardDeduction)

	for _, b := range c.Brackets {
		if taxableIncome.LessThanOrEqual(b.Lower) {
			break
		}

		tax = tax.Add(b.CalculateTax(taxableIncome))

		if taxableIncome.LessThan(b.Upper) {
			break
		}
	}
	return tax
}
