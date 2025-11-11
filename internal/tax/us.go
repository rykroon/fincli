package tax

import (
	"github.com/shopspring/decimal"
)

type UsFilingConfig struct {
	StandardDeduction decimal.Decimal `json:"standard_deduction"`
	Schedule          ProgressiveTax  `json:"schedule"`
}

// type UsTaxSystemResult struct {
// 	StandardDeduction   decimal.Decimal
// 	AdjustedGrossIncome decimal.Decimal
// 	TaxableIncome       decimal.Decimal
// 	MarginalTaxRate     decimal.Decimal
// 	TaxesDue            decimal.Decimal
// }

type UsTaxSystem struct {
	FilingConfigs map[FilingStatus]UsFilingConfig `json:"filing_configs"`
}

func (sys UsTaxSystem) CalculateTax(p TaxPayer) TaxResult {
	config, ok := sys.FilingConfigs[p.FilingStatus]
	if !ok {
		panic("Filing status not found.")
	}

	adjustedGrossIncome := p.Income
	for _, adj := range p.Adjustments {
		adjustedGrossIncome = adj.Adjust(adjustedGrossIncome)
	}

	taxableIncome := adjustedGrossIncome.Sub(config.StandardDeduction)
	marginalBracket := config.Schedule.GetMarginalBracket(taxableIncome)
	taxesDue := config.Schedule.CalculateTax(taxableIncome)

	result := NewTaxResult("Federal Tax", taxesDue)
	result.AddStat("Adjusted Gross Income", adjustedGrossIncome)
	result.AddStat("Standard Deduction", config.StandardDeduction)
	result.AddStat("Taxable Income", taxableIncome)
	result.AddStat("Marginal Tax Rate", marginalBracket.Rate)

	return result
}
