package tax

import "fmt"

// https://www.nj.gov/treasury/taxation/njit5.shtml

type NjTaxSystem struct {
	FilingConfigs map[FilingStatus]ProgressiveTax `json:"filing_configs"`
}

func (sys NjTaxSystem) CalculateTax(p TaxPayer) (TaxResult, error) {
	taxCalc, ok := sys.FilingConfigs[p.FilingStatus]
	if !ok {
		return TaxResult{}, fmt.Errorf(
			"no NJ tax data for filing status '%s'", p.FilingStatus,
		)
	}

	adjustedGrossIncome := p.Income
	for _, adj := range p.Adjustments {
		// to do: check for specific adjustments
		adjustedGrossIncome = adj.Adjust(adjustedGrossIncome)
	}

	taxesDue := taxCalc.CalculateTax(adjustedGrossIncome)
	result := NewTaxResult("NJ Tax", taxesDue)
	result.AddStat("Adj Gross Income", adjustedGrossIncome, "currency")
	result.AddStat("Marginal Tax Rate", taxCalc.GetMarginalBracket(adjustedGrossIncome).Rate, "percent")
	result.AddStat("Effective Tax Rate", taxesDue.Div(p.Income), "percent")
	result.AddStat("Taxes", taxesDue, "currency")
	return result, nil
}
