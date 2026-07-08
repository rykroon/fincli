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

type ExtraMonthlyStrategy struct {
	extraPayment decimal.Decimal
}

func NewExtraMonthlyStrategy(extraPayment decimal.Decimal) PaymentStrategy {
	return ExtraMonthlyStrategy{extraPayment: extraPayment}
}

func (s ExtraMonthlyStrategy) NextPayment(loan *Loan, sched *Schedule) decimal.Decimal {
	payment := CalculateMonthlyPayment(loan.Principal, loan.MonthlyRate(), loan.NumPeriods())
	payment = payment.Add(s.extraPayment)
	return payment
}

type ExtraAnnualStrategy struct {
	extraPayment decimal.Decimal
}

func NewExtraAnnualStrategy(extraPayment decimal.Decimal) PaymentStrategy {
	return ExtraAnnualStrategy{extraPayment: extraPayment}
}

func (s ExtraAnnualStrategy) NextPayment(loan *Loan, sched *Schedule) decimal.Decimal {
	payment := CalculateMonthlyPayment(loan.Principal, loan.MonthlyRate(), loan.NumPeriods())
	if len(sched.Payments)%12 == 0 {
		payment = payment.Add(s.extraPayment)
	}
	return payment
}
