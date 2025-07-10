package mortgage

import (
	"github.com/shopspring/decimal"
)

type Schedule struct {
	MonthlyPayment  decimal.Decimal
	StartingBalance decimal.Decimal
	Payments        []*Payment
	TotalAmount     decimal.Decimal
	TotalInterest   decimal.Decimal
}

func (s *Schedule) appendPayment(p *Payment) {
	s.Payments = append(s.Payments, p)
	s.TotalAmount = s.TotalAmount.Add(p.Total())
	s.TotalInterest = s.TotalInterest.Add(p.Interest())
}

func (s *Schedule) NumPeriods() decimal.Decimal {
	return decimal.NewFromInt(int64(len(s.Payments)))
}

func (s *Schedule) AverageMonthlyPayment() decimal.Decimal {
	return s.TotalAmount.Div(s.NumPeriods())
}

func CalculateSchedule(p decimal.Decimal, i decimal.Decimal, n decimal.Decimal, extraPaymentStratgey ExtraPaymentStrategy) Schedule {
	balance := p
	schedule := Schedule{
		MonthlyPayment:  CalculateMonthlyPayment(p, i, n),
		StartingBalance: balance,
	}
	for period := 1; balance.Round(2).GreaterThan(decimal.Zero); period++ {
		interest := balance.Mul(i)
		principal := schedule.MonthlyPayment.Sub(interest)
		payment := newPayment(period, principal, interest, balance)
		extraPrincipal := extraPaymentStratgey(period, principal, interest)
		payment.SetExtraPrincipal(extraPrincipal)

		balance = payment.Balance()
		schedule.appendPayment(payment)
	}
	return schedule
}
