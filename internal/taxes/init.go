package taxes

const maxIncome int64 = 1_000_000_000_000_000

func buildSingle2025() FilingConfig {
	return NewFilingConfig(
		15000,
		NewBracket(0, 11_925, .10),
		NewBracket(11_925, 48_475, .12),
		NewBracket(48_475, 103_350, .22),
		NewBracket(103_350, 197_300, .24),
		NewBracket(197_300, 250_525, .32),
		NewBracket(250_525, 626_350, .35),
		NewBracket(626_350, maxIncome, .37),
	)
}

func buildMarriedJointly2025() FilingConfig {
	return NewFilingConfig(
		31_500,
		NewBracket(0, 23_850, .10),
		NewBracket(23_850, 96_950, .12),
		NewBracket(96_950, 206_700, .22),
		NewBracket(206_700, 394_600, .24),
		NewBracket(394_600, 501_050, .32),
		NewBracket(501_050, 751_600, .35),
		NewBracket(751_600, maxIncome, .37),
	)
}

var UsFederalTaxTable = TaxTable{
	2025: {
		Single:       buildSingle2025(),
		MarriedJoint: buildMarriedJointly2025(),
	},
}

// There are some nuances to these taxes at higher incomes, so it's not actually a flat tax.
var SocialSecurityTax = NewBracket(0, 176_100, .062)
var MedicareTax = NewFlatTax(.0145)
