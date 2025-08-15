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

type FilingConfig struct {
	Brackets          []Bracket
	StandardDeduction decimal.Decimal
}

func (c FilingConfig) GetBracketByIncome(i decimal.Decimal) Bracket {
	for _, bracket := range c.Brackets {
		if bracket.Max.IsZero() {
			return bracket
		}
		if bracket.Min.LessThan(i) && i.LessThan(bracket.Max) {
			return bracket
		}
	}
	panic("a bracket should have been found")
}

func (c FilingConfig) CalculateTax(income decimal.Decimal) decimal.Decimal {
	tax := decimal.Zero
	income = income.Sub(c.StandardDeduction)

	for _, b := range c.Brackets {
		if income.LessThanOrEqual(b.Min) {
			break
		}

		upper := b.Max
		if upper.IsZero() || income.LessThan(upper) {
			upper = income
		}

		taxable := upper.Sub(b.Min)
		tax = tax.Add(taxable.Mul(b.Rate))

		if income.LessThan(upper) {
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

var UsFederalTaxTable = TaxTable{
	Years: map[int]TaxYear{
		2025: {
			Year: 2025,
			Filings: map[FilingStatus]FilingConfig{
				Single: {
					Brackets: []Bracket{
						{Min: decimal.Zero, Max: decimal.NewFromInt(11925), Rate: decimal.NewFromFloat(.10)},
						{Min: decimal.NewFromInt(11925), Max: decimal.NewFromInt(48475), Rate: decimal.NewFromFloat(.12)},
						{Min: decimal.NewFromInt(48475), Max: decimal.NewFromInt(103350), Rate: decimal.NewFromFloat(.22)},
						{Min: decimal.NewFromInt(103350), Max: decimal.NewFromInt(197300), Rate: decimal.NewFromFloat(.24)},
						{Min: decimal.NewFromInt(197300), Max: decimal.NewFromInt(250525), Rate: decimal.NewFromFloat(.32)},
						{Min: decimal.NewFromInt(197300), Max: decimal.NewFromInt(626350), Rate: decimal.NewFromFloat(.35)},
						{Min: decimal.NewFromInt(626350), Max: decimal.Zero, Rate: decimal.NewFromFloat(.37)},
					},
					StandardDeduction: decimal.NewFromInt(15000),
				},
			},
		},
	},
}
