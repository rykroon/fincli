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

type ExtraMonthlyPayment struct {
	extraPayment decimal.Decimal
}

func (s ExtraMonthlyPayment) NextPayment(loan *Loan) decimal.Decimal {
	payment := CalculateMonthlyPayment(loan.Principal, loan.MonthlyRate(), loan.NumPeriods())
	payment = payment.Add(s.extraPayment)
	return payment
}
