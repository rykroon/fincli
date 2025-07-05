package mortgage

import (
	"github.com/shopspring/decimal"
)

type Schedule struct {
	MonthlyPayment  decimal.Decimal
	StartingBalance decimal.Decimal
	Payments        []Payment
	TotalAmount     decimal.Decimal
	TotalInterest   decimal.Decimal
}

func (s *Schedule) appendPayment(p Payment) {
	s.Payments = append(s.Payments, p)
	s.TotalAmount = s.TotalAmount.Add(p.Total())
	s.TotalInterest = s.TotalInterest.Add(p.Interest)
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
	for period := 1; balance.GreaterThan(decimal.Zero); period++ {
		interest := balance.Mul(i)
		principal := schedule.MonthlyPayment.Sub(interest)

		extraPrincipal := extraPaymentStratgey(period, principal, interest)
		// make sure the principal doesn't go above the balance.
		if principal.Add(extraPrincipal).GreaterThan(balance) {
			if principal.GreaterThan(balance) {
				principal = balance
				extraPrincipal = decimal.Zero
			} else {
				extraPrincipal = balance.Sub(principal)
			}
		}

		balance = balance.Sub(principal.Add(extraPrincipal))
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
