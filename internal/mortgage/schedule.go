package mortgage

import (
	"github.com/shopspring/decimal"
)

type Schedule struct {
	Loan           *Loan
	Payments       []Payment
	TotalPrincipal decimal.Decimal
	TotalInterest  decimal.Decimal
}

func NewSchedule(loan *Loan) *Schedule {
	return &Schedule{
		Loan:     loan,
		Payments: make([]Payment, 0, loan.NumPeriods()),
	}
}

func (s *Schedule) appendPayment(p Payment) {
	s.Payments = append(s.Payments, p)
	s.TotalPrincipal = s.TotalPrincipal.Add(p.Principal)
	s.TotalInterest = s.TotalInterest.Add(p.Interest)
}

func (s *Schedule) NumPeriods() decimal.Decimal {
	return decimal.NewFromInt(int64(len(s.Payments)))
}

func (s *Schedule) TotalAmount() decimal.Decimal {
	return s.TotalPrincipal.Add(s.TotalInterest)
}

func (s *Schedule) RemainingBalance() decimal.Decimal {
	return s.Loan.Principal.Sub(s.TotalPrincipal)
}

func (s *Schedule) AverageMonthlyPayment() decimal.Decimal {
	return s.TotalAmount().Div(s.NumPeriods())
}

func CalculateSchedule(loan *Loan, strategy PaymentStrategy) *Schedule {
	schedule := NewSchedule(loan)
	balance := loan.Principal
	for period := 1; balance.Round(2).GreaterThan(decimal.Zero); period++ {
		paymentAmount := strategy.NextPayment(loan, schedule)
		interest := balance.Mul(loan.MonthlyRate())
		principal := paymentAmount.Sub(interest)
		balance = balance.Sub(principal)
		payment := NewPayment(period, principal, interest, balance)
		schedule.appendPayment(payment)
	}
	return schedule
}

func CalculateMonthlyPayment(p, i decimal.Decimal, n uint16) decimal.Decimal {
	// P * ((i * (1+i)^n) / ((1+i)^n - 1))
	one := decimal.NewFromInt(1)
	onePlusIPowN := one.Add(i).Pow(decimal.NewFromUint64(uint64(n)))
	return p.Mul((i.Mul(onePlusIPowN)).Div(onePlusIPowN.Sub(one)))
}
