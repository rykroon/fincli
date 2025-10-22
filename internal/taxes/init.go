package taxes

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
	2025: {
		Single:       buildSingle2025(),
		MarriedJoint: buildMarriedJointly2025(),
	},
}
