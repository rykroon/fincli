package mortgage

import (
	"github.com/shopspring/decimal"
)

type Loan struct {
	Principal  decimal.Decimal
	AnnualRate decimal.Decimal
	NumYears   int64
}

func NewLoan(p, r decimal.Decimal, y int64) Loan {
	return Loan{
		Principal:  p,
		AnnualRate: r,
		NumYears:   y,
	}
}

func (l Loan) MonthlyRate() decimal.Decimal {
	return l.AnnualRate.Div(decimal.NewFromInt(12))
}

func (l Loan) NumPeriods() int64 {
	return l.NumYears * 12
}
