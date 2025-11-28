package mortgage

import (
	"github.com/shopspring/decimal"
)

type PaymentStrategy interface {
	NextPayment(*Loan, *Schedule) decimal.Decimal
}

type DefaultStrategy struct{}

func NewDefaultStrategy() PaymentStrategy {
	return DefaultStrategy{}
}

func (s DefaultStrategy) NextPayment(loan *Loan, sched *Schedule) decimal.Decimal {
	return CalculateMonthlyPayment(loan.Principal, loan.MonthlyRate(), loan.NumPeriods())
}

type ExtraMonthlyStratgey struct {
	extraPayment decimal.Decimal
}

func NewExtraMonthlyStrategy(extraPayment decimal.Decimal) PaymentStrategy {
	return ExtraMonthlyStratgey{extraPayment: extraPayment}
}

func (s ExtraMonthlyStratgey) NextPayment(loan *Loan, sched *Schedule) decimal.Decimal {
	payment := CalculateMonthlyPayment(loan.Principal, loan.MonthlyRate(), loan.NumPeriods())
	payment = payment.Add(s.extraPayment)
	return payment
}

type ExtraAnnualStratgey struct {
	extraPayment decimal.Decimal
}

func NewExtraAnnualStrategy(extraPayment decimal.Decimal) PaymentStrategy {
	return ExtraAnnualStratgey{extraPayment: extraPayment}
}

func (s ExtraAnnualStratgey) NextPayment(loan *Loan, sched *Schedule) decimal.Decimal {
	payment := CalculateMonthlyPayment(loan.Principal, loan.MonthlyRate(), loan.NumPeriods())
	if len(sched.Payments)%12 == 0 {
		payment = payment.Add(s.extraPayment)
	}
	return payment
}
