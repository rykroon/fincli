package tax

import (
	"github.com/shopspring/decimal"
)

const maxIncome int = 1_000_000_000_000_000

func buildUsTaxSystem2025() UsTaxSystem {
	sys := NewUsTaxSystem()
	sys.AddFilingStatus(Single, decimal.NewFromInt(15_000), buildUsSingle2025())
	sys.AddFilingStatus(MarriedJoint, decimal.NewFromInt(31_500), buildUsMarriedJointly2025())
	return sys
}

func buildUsSingle2025() ProgressiveTax {
	t := NewProgressiveTax()
	t.AddBracket(newBracket(0, 11_925, .10))
	t.AddBracket(newBracket(11_925, 48_475, .12))
	t.AddBracket(newBracket(48_475, 103_350, .22))
	t.AddBracket(newBracket(103_350, 197_300, .24))
	t.AddBracket(newBracket(197_300, 250_525, .32))
	t.AddBracket(newBracket(250_525, 626_350, .35))
	t.AddBracket(newBracket(626_350, maxIncome, .37))
	return t
}

func buildUsMarriedJointly2025() ProgressiveTax {
	t := NewProgressiveTax()
	t.AddBracket(newBracket(0, 23_850, .10))
	t.AddBracket(newBracket(23_850, 96_950, .12))
	t.AddBracket(newBracket(96_950, 206_700, .22))
	t.AddBracket(newBracket(206_700, 394_600, .24))
	t.AddBracket(newBracket(394_600, 501_050, .32))
	t.AddBracket(newBracket(501_050, 751_600, .35))
	t.AddBracket(newBracket(751_600, maxIncome, .37))
	return t
}

var UsFederalRegistry = map[uint16]UsTaxSystem{
	2025: buildUsTaxSystem2025(),
}

var FicaRegistry = map[uint16]FicaTaxSystem{
	2025: buildFica2025(),
}

func buildFica2025() FicaTaxSystem {
	return FicaTaxSystem{
		SocialSecurityTax: NewCappedTax(.062, 176_100),
		MedicareTax:       NewFlatTax(.0145),
	}
}
