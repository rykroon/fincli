package tax

// https://www.nj.gov/treasury/taxation/njit5.shtml

type NjTaxSystem struct {
	FilingConfigs map[FilingStatus]ProgressiveTax `json:"filing_configs"`
}

func (sys NjTaxSystem) CalculateTax(p TaxPayer) TaxResult {
	taxCalc, ok := sys.FilingConfigs[p.FilingStatus]
	if !ok {
		panic("filing status not found")
	}

	adjustedGrossIncome := p.Income
	for _, adj := range p.Adjustments {
		// to do: check for specific adjustments
		adjustedGrossIncome = adj.Adjust(adjustedGrossIncome)
	}

	taxesDue := taxCalc.CalculateTax(adjustedGrossIncome)
	result := NewTaxResult("NJ Tax", taxesDue)
	result.AddStat("Adj Gross Income", adjustedGrossIncome, "currency")
	result.AddStat("Marginal Tax Rate", taxCalc.GetMarginalBracket(p.Income).Rate, "percent")
	result.AddStat("Effective Tax Rate", taxesDue.Div(p.Income), "percent")
	result.AddStat("Taxes", taxesDue, "currency")
	return result
}
