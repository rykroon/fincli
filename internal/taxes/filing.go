package taxes

import "github.com/shopspring/decimal"

type FilingConfig struct {
	StandardDeduction decimal.Decimal
	Calculator        TaxCalculator
}

func NewFilingConfig(standardDecution float64, calc TaxCalculator) FilingConfig {
	return FilingConfig{decimal.NewFromFloat(standardDecution), calc}
}

func (c FilingConfig) CalculateTax(income decimal.Decimal) decimal.Decimal {
	taxableIncome := income.Sub(c.StandardDeduction)
	return c.Calculator.CalculateTax(taxableIncome)
}
