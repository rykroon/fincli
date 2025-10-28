package tax

import (
	"github.com/shopspring/decimal"
)

type UsFilingConfig struct {
	StandardDeduction decimal.Decimal `json:"standard_deduction"`
	Schedule          ProgressiveTax  `json:"schedule"`
}

func NewUsFilingConfig(standardDeduction int64, sched ProgressiveTax) UsFilingConfig {
	return UsFilingConfig{
		StandardDeduction: decimal.NewFromInt(standardDeduction),
		Schedule:          sched,
	}
}

type UsTaxSystemResult struct {
	StandardDeduction   decimal.Decimal
	AdjustedGrossIncome decimal.Decimal
	TaxableIncome       decimal.Decimal
	MarginalTaxRate     decimal.Decimal
	TaxesDue            decimal.Decimal
}

type UsTaxSystem struct {
	FilingConfigs map[FilingStatus]UsFilingConfig `json:"filing_configs"`
}

func NewUsTaxSystem() UsTaxSystem {
	return UsTaxSystem{
		FilingConfigs: make(map[FilingStatus]UsFilingConfig),
	}
}

func (sys *UsTaxSystem) AddFilingStatus(status FilingStatus, config UsFilingConfig) *UsTaxSystem {
	sys.FilingConfigs[status] = config
	return sys
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
	marginalBracket := config.Schedule.GetMarginalBracket(taxableIncome)
	taxesDue := config.Schedule.CalculateTax(taxableIncome)

	return UsTaxSystemResult{
		StandardDeduction:   config.StandardDeduction,
		AdjustedGrossIncome: adjustedGrossIncome,
		TaxableIncome:       taxableIncome,
		MarginalTaxRate:     marginalBracket.Rate,
		TaxesDue:            taxesDue,
	}
}
