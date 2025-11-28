package mortgage

import (
	"github.com/shopspring/decimal"
)

type Loan struct {
	Principal  decimal.Decimal
	AnnualRate decimal.Decimal
	NumYears   uint16
}

func NewLoan(p, r decimal.Decimal, y uint16) *Loan {
	return &Loan{
		Principal:  p,
		AnnualRate: r,
		NumYears:   y,
	}
}

func (l Loan) MonthlyRate() decimal.Decimal {
	return l.AnnualRate.Div(decimal.NewFromInt(12))
}

func (l Loan) NumPeriods() uint16 {
	return l.NumYears * 12
}
