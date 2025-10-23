package taxes

import "github.com/shopspring/decimal"

type TaxCalculator interface {
	CalculateTax(decimal.Decimal) decimal.Decimal
}

type TaxSchedule map[FilingStatus]TaxCalculator

type TaxRegistry map[int]TaxSchedule
