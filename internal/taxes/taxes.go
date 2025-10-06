package taxes

import "github.com/shopspring/decimal"

type FilingStatus string

const (
	Single          FilingStatus = "single"
	MarriedJoint    FilingStatus = "married_joint"
	MarriedSeparate FilingStatus = "married_separate"
	HeadOfHouse     FilingStatus = "head_of_household"
)

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

type TaxYear struct {
	Year    int
	Filings map[FilingStatus]FilingConfig
}

type TaxTable struct {
	Years map[int]TaxYear
}

func (t *TaxTable) AddTaxYear(ty TaxYear) {
	t.Years[ty.Year] = ty
}

func (t TaxTable) GetConfig(year int, status FilingStatus) (*FilingConfig, bool) {
	taxYear, ok := t.Years[year]
	if !ok {
		return nil, false
	}
	config, ok := taxYear.Filings[status]
	if !ok {
		return nil, false
	}
	return &config, true
}

func buildSingle2025() FilingConfig {
	fc := FilingConfig{}
	fc.setStandardDeductionInt(15000)
	fc.addBracket(NewBracket(0, 11925, .10))
	fc.addBracket(NewBracket(11925, 48475, .12))
	fc.addBracket(NewBracket(48475, 103350, .22))
	fc.addBracket(NewBracket(103350, 197300, .24))
	fc.addBracket(NewBracket(197300, 250525, .32))
	fc.addBracket(NewBracket(250525, 626350, .35))
	fc.addBracket(NewBracket(626350, 0, .37))
	return fc
}

func buildMarriedJointly2025() FilingConfig {
	fc := FilingConfig{}
	fc.setStandardDeductionInt(31_500)
	fc.addBracket(NewBracket(0, 23_850, .10))
	fc.addBracket(NewBracket(23_850, 96_950, .12))
	fc.addBracket(NewBracket(96_950, 206_700, .22))
	fc.addBracket(NewBracket(206_700, 394_600, .24))
	fc.addBracket(NewBracket(394_600, 501_050, .32))
	fc.addBracket(NewBracket(501_050, 751_600, .35))
	fc.addBracket(NewBracket(751_600, 0, .37))
	return fc
}

var UsFederalTaxTable = TaxTable{
	Years: map[int]TaxYear{
		2025: {
			Year: 2025,
			Filings: map[FilingStatus]FilingConfig{
				Single:       buildSingle2025(),
				MarriedJoint: buildMarriedJointly2025(),
			},
		},
	},
}
