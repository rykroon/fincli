package tax

import (
	"github.com/rykroon/fincli/internal/tax/calculators"
	"github.com/shopspring/decimal"
)

type FicaTaxSystem struct {
	SocialSecurityTax calculators.CappedTax
	MedicareTax       calculators.FlatTax
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
