package tax

import "github.com/shopspring/decimal"

type Adjustment struct {
	Label  string
	Amount decimal.Decimal
}

func (a Adjustment) Adjust(income decimal.Decimal) decimal.Decimal {
	return income.Sub(a.Amount)
}
