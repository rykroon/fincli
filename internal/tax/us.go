package tax

import (
	"github.com/shopspring/decimal"
)

type UsTaxSystem struct {
	StandardDeductions map[FilingStatus]decimal.Decimal
	Calculators        map[FilingStatus]ProgressiveTax
}

type UsTaxSystemResult struct {
	StandardDeduction   decimal.Decimal
	AdjustedGrossIncome decimal.Decimal
	TaxableIncome       decimal.Decimal
	MarginalTaxRate     decimal.Decimal
	TaxesDue            decimal.Decimal
}

func (s UsTaxSystem) CalculateTax(p TaxPayer) UsTaxSystemResult {
	adjustedGrossIncome := p.Income
	for _, adj := range p.Adjustments {
		adjustedGrossIncome = adj.Adjust(adjustedGrossIncome)
	}

	standardDeduction, ok := s.StandardDeductions[p.FilingStatus]
	if !ok {
		panic("Filing status not found.")
	}
	taxableIncome := adjustedGrossIncome.Sub(standardDeduction)

	calc, ok := s.Calculators[p.FilingStatus]
	if !ok {
		panic("Filing status not found.")
	}

	marginalBracket := calc.GetMarginalBracket(taxableIncome)
	taxesDue := calc.CalculateTax(taxableIncome)
	return UsTaxSystemResult{
		StandardDeduction:   standardDeduction,
		AdjustedGrossIncome: adjustedGrossIncome,
		TaxableIncome:       taxableIncome,
		MarginalTaxRate:     marginalBracket.Rate,
		TaxesDue:            taxesDue,
	}
}
