package tax

import "github.com/shopspring/decimal"

type NjTaxSystem struct {
	FilingConfigs map[FilingStatus]ProgressiveTax `json:"filing_configs"`
}

type NjTaxResult struct {
	MarginalTaxRate decimal.Decimal
	TaxesDue        decimal.Decimal
}

func (sys NjTaxSystem) CalculateTax(p TaxPayer) TaxResult {
	taxCalc, ok := sys.FilingConfigs[p.FilingStatus]
	if !ok {
		panic("filing status not found")
	}
	taxesDue := taxCalc.CalculateTax(p.Income)
	result := NewTaxResult("NJ Tax", taxesDue)
	result.AddStat("Marginal Tax Rate", taxCalc.GetMarginalBracket(p.Income).Rate)
	return result
}
