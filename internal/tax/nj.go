package tax

import "github.com/shopspring/decimal"

type NjTaxSystem struct {
	FilingConfigs map[FilingStatus]ProgressiveTax `json:"filing_configs"`
}

type NjTaxResult struct {
	MarginalTaxRate decimal.Decimal
	TaxesDue        decimal.Decimal
}

func (sys NjTaxSystem) CalculateTax(p TaxPayer) NjTaxResult {
	taxCalc, ok := sys.FilingConfigs[p.FilingStatus]
	if !ok {
		panic("filing status not found")
	}
	return NjTaxResult{
		MarginalTaxRate: taxCalc.GetMarginalBracket(p.Income).Rate,
		TaxesDue:        taxCalc.CalculateTax(p.Income),
	}
}
