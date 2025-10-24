package tax

import "github.com/shopspring/decimal"

type FilingStatus string

const (
	Single          FilingStatus = "single"
	MarriedJoint    FilingStatus = "married_joint"
	MarriedSeparate FilingStatus = "married_separate"
	HeadOfHouse     FilingStatus = "head_of_household"
)

type TaxPayer struct {
	Income       decimal.Decimal
	FilingStatus FilingStatus
	Adjustments  []Adjustment
}

func NewTaxPayer(income decimal.Decimal, filingStatus FilingStatus, adjustments ...Adjustment) TaxPayer {
	return TaxPayer{income, filingStatus, adjustments}
}
