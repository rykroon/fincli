package tax

type FicaTaxSystem struct {
	SocialSecurityTax CappedTax `json:"social_security_tax"`
	MedicareTax       FlatTax   `json:"medicare_tax"`
}

func (s FicaTaxSystem) CalculateTax(p TaxPayer) TaxResult {
	ssTaxDue := s.SocialSecurityTax.CalculateTax(p.Income)
	mcTaxDue := s.MedicareTax.CalculateTax(p.Income)

	total := ssTaxDue.Add(mcTaxDue)

	result := NewTaxResult("FICA Tax", total)
	result.AddStat("Social Security Tax", ssTaxDue, "currency")
	result.AddStat("Medicare Tax", mcTaxDue, "currency")
	return result
}
