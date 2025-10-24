package tax

import (
	"github.com/rykroon/fincli/internal/tax/calculators"
	"github.com/shopspring/decimal"
)

type TaxSystem interface {
	CalculateTax(TaxPayer) decimal.Decimal
}

type UsTaxSystem struct {
	StandardDeductions map[FilingStatus]decimal.Decimal
	Calculators        map[FilingStatus]calculators.TaxCalculator
}

func (s UsTaxSystem) CalculateTax(p TaxPayer) decimal.Decimal {
	taxableIncome := p.Income
	for _, adj := range p.Adjustments {
		taxableIncome = adj.Adjust(taxableIncome)
	}

	standardDeduction, ok := s.StandardDeductions[p.FilingStatus]
	if !ok {
		panic("Filing status not found.")
	}
	taxableIncome = taxableIncome.Sub(standardDeduction)

	calc, ok := s.Calculators[p.FilingStatus]
	if !ok {
		panic("Filing status not found.")
	}
	return calc.CalculateTax(taxableIncome)
}
