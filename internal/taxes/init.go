package taxes

const maxIncome int64 = 1_000_000_000_000_000

func buildSingle2025() FilingConfig {
	fc := FilingConfig{}
	fc.setStandardDeductionInt(15000)
	fc.addBracket(NewBracket(0, 11_925, .10))
	fc.addBracket(NewBracket(11_925, 48_475, .12))
	fc.addBracket(NewBracket(48_475, 103_350, .22))
	fc.addBracket(NewBracket(103_350, 197_300, .24))
	fc.addBracket(NewBracket(197_300, 250_525, .32))
	fc.addBracket(NewBracket(250_525, 626_350, .35))
	fc.addBracket(NewBracket(626_350, maxIncome, .37))
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
	fc.addBracket(NewBracket(751_600, maxIncome, .37))
	return fc
}

var UsFederalTaxTable = TaxTable{
	2025: {
		Single:       buildSingle2025(),
		MarriedJoint: buildMarriedJointly2025(),
	},
}
