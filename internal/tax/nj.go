package tax

import "github.com/shopspring/decimal"

type NjTaxSystem struct {
	FilingConfigs map[FilingStatus]ProgressiveTax `json:"filing_configs"`
}

func NewNjTaxSystem() NjTaxSystem {
	return NjTaxSystem{FilingConfigs: make(map[FilingStatus]ProgressiveTax)}
}

func (sys *NjTaxSystem) AddFilingConfig(status FilingStatus, tax ProgressiveTax) {
	sys.FilingConfigs[status] = tax
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
