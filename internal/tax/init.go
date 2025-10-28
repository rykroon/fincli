package tax

const maxIncome uint64 = 1_000_000_000_000_000

func buildUsTaxSystem2025() UsTaxSystem {
	sys := NewUsTaxSystem()
	sys.AddFilingStatus(Single, NewUsFilingConfig(15_000, buildUsSingle2025()))
	sys.AddFilingStatus(MarriedJoint, NewUsFilingConfig(31_500, buildUsMarriedJointly2025()))
	return sys
}

func buildUsSingle2025() ProgressiveTax {
	return NewProgressiveTax(
		NewBracket(0, 11_925, .10),
		NewBracket(11_925, 48_475, .12),
		NewBracket(48_475, 103_350, .22),
		NewBracket(103_350, 197_300, .24),
		NewBracket(197_300, 250_525, .32),
		NewBracket(250_525, 626_350, .35),
		NewBracket(626_350, maxIncome, .37),
	)
}

func buildUsMarriedJointly2025() ProgressiveTax {
	return NewProgressiveTax(
		NewBracket(0, 23_850, .10),
		NewBracket(23_850, 96_950, .12),
		NewBracket(96_950, 206_700, .22),
		NewBracket(206_700, 394_600, .24),
		NewBracket(394_600, 501_050, .32),
		NewBracket(501_050, 751_600, .35),
		NewBracket(751_600, maxIncome, .37),
	)
}

func buildFica2025() FicaTaxSystem {
	return FicaTaxSystem{
		SocialSecurityTax: NewCappedTax(176_100, .062),
		MedicareTax:       NewFlatTax(.0145),
	}
}

func buildNewJerseyTaxSystem() NjTaxSystem {
	sys := NewNjTaxSystem()
	sys.AddFilingConfig(Single, buildNjSingle2025())
	sys.AddFilingConfig(MarriedSeparate, buildNjSingle2025())
	sys.AddFilingConfig(MarriedJoint, buildNjMarried2025())
	sys.AddFilingConfig(HeadOfHouse, buildNjMarried2025())
	return sys
}

func buildNjSingle2025() ProgressiveTax {
	return NewProgressiveTax(
		NewBracket(0, 20_000, .014),
		NewBracket(20_000, 35_000, .0175),
		NewBracket(35_000, 40_000, .035),
		NewBracket(40_000, 75_000, .05525),
		NewBracket(75_000, 500_000, .0637),
		NewBracket(500_000, 1_000_000, .0897),
		NewBracket(1_000_000, maxIncome, .1075),
	)
}

func buildNjMarried2025() ProgressiveTax {
	return NewProgressiveTax(
		NewBracket(0, 20_000, .014),
		NewBracket(20_000, 50_000, .0175),
		NewBracket(50_000, 70_000, .0245),
		NewBracket(70_000, 80_000, .035),
		NewBracket(80_000, 150_000, .05525),
		NewBracket(150_000, 500_000, .0637),
		NewBracket(500_000, 1_000_000, .0897),
		NewBracket(1_000_000, maxIncome, .1075),
	)
}

var UsFederalRegistry = map[uint16]UsTaxSystem{
	2025: buildUsTaxSystem2025(),
}

var FicaRegistry = map[uint16]FicaTaxSystem{
	2025: buildFica2025(),
}

var NewJerseyRegistry = map[uint16]NjTaxSystem{
	2025: buildNewJerseyTaxSystem(),
}
