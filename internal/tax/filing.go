package tax

// type FilingConfig struct {
// 	StandardDeduction decimal.Decimal
// 	Calculator        calculators.TaxCalculator
// }

// func NewFilingConfig(standardDecution float64, calc calculators.TaxCalculator) FilingConfig {
// 	return FilingConfig{decimal.NewFromFloat(standardDecution), calc}
// }

// func (c FilingConfig) CalculateTax(income decimal.Decimal) decimal.Decimal {
// 	taxableIncome := income.Sub(c.StandardDeduction)
// 	return c.Calculator.CalculateTax(taxableIncome)
// }

/*
 * considering the idea that a FilingConfig is not a Calculator itself,
 * but can provide a calculator if you give it additional configuration information.
 * For example someone can choose to itemize instead of using the standard deduction.
 * Addtionally there are 401k contributions HSA contributions and other things that reduce
 * someone's taxable income.
 *
 */
