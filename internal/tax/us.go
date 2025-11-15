package tax

import (
	"github.com/shopspring/decimal"
)

// https://www.eitc.irs.gov/publications/p17#d0e50262

type UsFilingConfig struct {
	StandardDeduction decimal.Decimal `json:"standard_deduction"`
	Schedule          ProgressiveTax  `json:"schedule"`
}

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
	result.AddStat("Adj Gross Income", adjustedGrossIncome, "currency")
	result.AddStat("Standard Deduction", config.StandardDeduction, "currency")
	result.AddStat("Taxable Income", taxableIncome, "currency")
	result.AddStat("Marginal Tax Rate", marginalBracket.Rate, "percent")
	result.AddStat("Effective Tax Rate", taxesDue.Div(p.Income), "percent")
	result.AddStat("Taxes", taxesDue, "currency")

	return result
}
