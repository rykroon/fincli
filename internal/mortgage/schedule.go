package mortgage

import (
	"math"
)

type Schedule struct {
	MonthlyPayment  float64
	StartingBalance float64
	Payments        []Payment
	TotalAmount     float64
	TotalInterest   float64
}

func (s *Schedule) appendPayment(p Payment) {
	s.Payments = append(s.Payments, p)
	s.TotalAmount += p.Total()
	s.TotalInterest += p.Interest
}

func (s *Schedule) NumPeriods() int {
	return len(s.Payments)
}

func (s *Schedule) AverageMonthlyPayment() float64 {
	return s.TotalAmount / float64(s.NumPeriods())
}

func CalculateSchedule(p float64, i float64, n int, extraPaymentStratgey ExtraPaymentStrategy) Schedule {
	balance := p
	schedule := Schedule{
		MonthlyPayment:  CalculateMonthlyPayment(p, i, n),
		StartingBalance: balance,
	}
	for period := 1; math.Round(balance) > 0; period++ {
		interest := balance * i
		principal := schedule.MonthlyPayment - interest

		extraPrincipal := extraPaymentStratgey(period, principal, interest)
		// make sure the principal doesn't go above the balance.
		if principal+extraPrincipal > balance {
			if principal > balance {
				principal = balance
				extraPrincipal = 0
			} else {
				extraPrincipal = balance - principal
			}
		}

		balance -= principal + extraPrincipal
		p := Payment{
			Period:         period,
			Principal:      principal,
			ExtraPrincipal: extraPrincipal,
			Interest:       interest,
			Balance:        balance,
		}
		schedule.appendPayment(p)
	}
	return schedule
}
