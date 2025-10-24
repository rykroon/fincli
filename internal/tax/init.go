package tax

import (
	"github.com/rykroon/fincli/internal/tax/calculators"
	"github.com/shopspring/decimal"
)

const maxIncome int64 = 1_000_000_000_000_000

func buildUsTaxSystem2025() UsTaxSystem {
	system := UsTaxSystem{
		map[FilingStatus]decimal.Decimal{
			Single:       decimal.NewFromInt(15_000),
			MarriedJoint: decimal.NewFromInt(31_500),
		},
		map[FilingStatus]calculators.TaxCalculator{
			Single:       buildUsSingle2025(),
			MarriedJoint: buildUsMarriedJointly2025(),
		},
	}
	return system
}

func buildUsSingle2025() calculators.ProgressiveTax {
	t := calculators.NewProgressiveTax()
	t.AddBracket(0, 11_925, .10)
	t.AddBracket(11_925, 48_475, .12)
	t.AddBracket(48_475, 103_350, .22)
	t.AddBracket(103_350, 197_300, .24)
	t.AddBracket(197_300, 250_525, .32)
	t.AddBracket(250_525, 626_350, .35)
	t.AddBracket(626_350, maxIncome, .37)
	return t
}

func buildUsMarriedJointly2025() calculators.ProgressiveTax {
	t := calculators.NewProgressiveTax()
	t.AddBracket(0, 23_850, .10)
	t.AddBracket(23_850, 96_950, .12)
	t.AddBracket(96_950, 206_700, .22)
	t.AddBracket(206_700, 394_600, .24)
	t.AddBracket(394_600, 501_050, .32)
	t.AddBracket(501_050, 751_600, .35)
	t.AddBracket(751_600, maxIncome, .37)
	return t
}

var UsFederalRegistry = TaxRegistry[UsTaxSystem]{
	2025: buildUsTaxSystem2025(),
}

// There are some nuances to these taxes at higher incomes, so it's not actually a flat tax.
var SocialSecurityTax = calculators.NewCappedTax(.062, 176_100)
var MedicareTax = calculators.NewFlatTax(.0145)
