package calculators

import "github.com/shopspring/decimal"

type TaxCalculator interface {
	CalculateTax(decimal.Decimal) decimal.Decimal
}
