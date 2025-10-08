package mortgage

import (
	"github.com/shopspring/decimal"
)

type Schedule struct {
	Loan          *Loan
	Payments      []Payment
	TotalAmount   decimal.Decimal
	TotalInterest decimal.Decimal
}

func (s *Schedule) appendPayment(p Payment) {
	s.Payments = append(s.Payments, p)
	s.TotalAmount = s.TotalAmount.Add(p.Total())
	s.TotalInterest = s.TotalInterest.Add(p.Interest)
}

func (s Schedule) NumPeriods() decimal.Decimal {
	return decimal.NewFromInt(int64(len(s.Payments)))
}

func (s Schedule) AverageMonthlyPayment() decimal.Decimal {
	return s.TotalAmount.Div(s.NumPeriods())
}

func CalculateSchedule(loan *Loan) Schedule {
	balance := loan.Principal
	schedule := Schedule{
		Loan:     loan,
		Payments: make([]Payment, 0, loan.NumPeriods()),
	}
	for period := 1; balance.Round(2).GreaterThan(decimal.Zero); period++ {
		paymentAmount := DefaultStrategy{}.NextPayment(loan)
		interest := balance.Mul(loan.MonthlyRate())
		principal := paymentAmount.Sub(interest)
		balance = balance.Sub(principal)
		payment := newPayment(period, principal, interest, balance)
		schedule.appendPayment(payment)
	}
	return schedule
}

func CalculateMonthlyPayment(p decimal.Decimal, i decimal.Decimal, n int64) decimal.Decimal {
	// P * ((i * (1+i)^n) / ((1+i)^n - 1))
	one := decimal.NewFromInt(1)
	onePlusIPowN := one.Add(i).Pow(decimal.NewFromInt(n))
	return p.Mul((i.Mul(onePlusIPowN)).Div(onePlusIPowN.Sub(one)))
}
