package tax

import (
	"github.com/shopspring/decimal"
)

type usFilingConfig struct {
	StandardDeduction decimal.Decimal
	Calculator        ProgressiveTax
}

type UsTaxSystem struct {
	FilingConfigs map[FilingStatus]usFilingConfig
}

func NewUsTaxSystem() UsTaxSystem {
	return UsTaxSystem{
		FilingConfigs: make(map[FilingStatus]usFilingConfig),
	}
}

type UsTaxSystemResult struct {
	StandardDeduction   decimal.Decimal
	AdjustedGrossIncome decimal.Decimal
	TaxableIncome       decimal.Decimal
	MarginalTaxRate     decimal.Decimal
	TaxesDue            decimal.Decimal
}

func (sys *UsTaxSystem) AddFilingStatus(status FilingStatus, standardDeduction decimal.Decimal, calc ProgressiveTax) {
	sys.FilingConfigs[status] = usFilingConfig{
		StandardDeduction: standardDeduction,
		Calculator:        calc,
	}
}

func (sys UsTaxSystem) CalculateTax(p TaxPayer) UsTaxSystemResult {
	adjustedGrossIncome := p.Income
	for _, adj := range p.Adjustments {
		adjustedGrossIncome = adj.Adjust(adjustedGrossIncome)
	}

	config, ok := sys.FilingConfigs[p.FilingStatus]
	if !ok {
		panic("Filing status not found.")
	}

	taxableIncome := adjustedGrossIncome.Sub(config.StandardDeduction)

	marginalBracket := config.Calculator.GetMarginalBracket(taxableIncome)
	taxesDue := config.Calculator.CalculateTax(taxableIncome)
	return UsTaxSystemResult{
		StandardDeduction:   config.StandardDeduction,
		AdjustedGrossIncome: adjustedGrossIncome,
		TaxableIncome:       taxableIncome,
		MarginalTaxRate:     marginalBracket.Rate,
		TaxesDue:            taxesDue,
	}
}
