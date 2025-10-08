package mortgage

import (
	"github.com/shopspring/decimal"
)

type PaymentStrategy interface {
	NextPayment(loan *Loan) decimal.Decimal
}

type DefaultStrategy struct{}

func (s DefaultStrategy) NextPayment(loan *Loan) decimal.Decimal {
	return CalculateMonthlyPayment(loan.Principal, loan.MonthlyRate(), loan.NumPeriods())
}
