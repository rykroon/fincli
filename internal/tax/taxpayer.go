package tax

import (
	"fmt"

	"github.com/shopspring/decimal"
)

type FilingStatus string

const (
	Single          FilingStatus = "single"
	MarriedJoint    FilingStatus = "married_joint"
	MarriedSeparate FilingStatus = "married_separate"
	HeadOfHouse     FilingStatus = "head_of_household"
)

func ParseFilingStatus(s string) (FilingStatus, error) {
	switch fs := FilingStatus(s); fs {
	case Single, MarriedJoint, MarriedSeparate, HeadOfHouse:
		return fs, nil
	default:
		return "", fmt.Errorf(
			"invalid filing status '%s', must be one of: %s, %s, %s, %s",
			s, Single, MarriedJoint, MarriedSeparate, HeadOfHouse,
		)
	}
}

type TaxPayer struct {
	Income       decimal.Decimal
	FilingStatus FilingStatus
	Adjustments  []Adjustment
}

func NewTaxPayer(income decimal.Decimal, filingStatus FilingStatus, adjustments ...Adjustment) TaxPayer {
	return TaxPayer{income, filingStatus, adjustments}
}
