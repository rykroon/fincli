package tax

import (
	"github.com/shopspring/decimal"
)

type FicaTaxSystem struct {
	SocialSecurityTax CappedTax `json:"social_security_tax"`
	MedicareTax       FlatTax   `json:"medicare_tax"`
}

type FicaTaxResult struct {
	SocialSecurityTaxDue decimal.Decimal
	MedicareTaxDue       decimal.Decimal
}

func (s FicaTaxSystem) CalculateTax(p TaxPayer) FicaTaxResult {
	return FicaTaxResult{
		SocialSecurityTaxDue: s.SocialSecurityTax.CalculateTax(p.Income),
		MedicareTaxDue:       s.MedicareTax.CalculateTax(p.Income),
	}
}
