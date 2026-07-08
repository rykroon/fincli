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

type ExtraPaymentStrategy struct {
	extraMonthly decimal.Decimal
	extraAnnual  decimal.Decimal
}

func NewExtraPaymentStrategy(extraMonthly, extraAnnual decimal.Decimal) PaymentStrategy {
	return ExtraPaymentStrategy{extraMonthly: extraMonthly, extraAnnual: extraAnnual}
}

func (s ExtraPaymentStrategy) NextPayment(loan *Loan, sched *Schedule) decimal.Decimal {
	payment := CalculateMonthlyPayment(loan.Principal, loan.MonthlyRate(), loan.NumPeriods())
	payment = payment.Add(s.extraMonthly)
	// the extra annual payment lands on the first month of each loan year
	if len(sched.Payments)%12 == 0 {
		payment = payment.Add(s.extraAnnual)
	}
	return payment
}
