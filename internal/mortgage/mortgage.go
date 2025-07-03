package mortgage

import (
	"math"
)

func CalculateMonthlyPayment(p float64, i float64, n int) float64 {
	return p * ((i * math.Pow(1+i, float64(n))) / (math.Pow(1+i, float64(n)) - 1))
}

type Payment struct {
	Period         int
	Principal      float64
	ExtraPrincipal float64
	Interest       float64
	Balance        float64
}

func (p *Payment) TotalPrincipal() float64 {
	return p.Principal + p.ExtraPrincipal
}

func (p *Payment) Total() float64 {
	return p.TotalPrincipal() + p.Interest
}

type paymentScheduleStats struct {
	TotalPayments         float64
	TotalInterest         float64
	AverageMonthlyPayment float64
}

func GetPaymentScheduleStats(payments []Payment) paymentScheduleStats {
	pss := paymentScheduleStats{}

	for _, p := range payments {
		pss.TotalInterest += p.Interest
		pss.TotalPayments += p.Total()
	}
	pss.AverageMonthlyPayment = pss.TotalPayments / float64(len(payments))
	return pss
}

func CalculatePayments(p float64, i float64, n int, extraPaymentStratgey ExtraPaymentStrategy) (float64, []Payment) {
	monthlyPayment := CalculateMonthlyPayment(p, i, n)
	payments := []Payment{}
	balance := p
	for period := 1; math.Round(balance) > 0; period++ {
		interest := balance * i
		principal := monthlyPayment - interest

		extraPrincipal := extraPaymentStratgey(period, principal, interest)
		if principal+extraPrincipal > balance {
			// make sure the extra payment does not pay for more balance there is.
			extraPrincipal = balance - principal
		}

		balance -= principal + extraPrincipal
		payments = append(payments, Payment{
			Period:         period,
			Principal:      principal,
			ExtraPrincipal: extraPrincipal,
			Interest:       interest,
			Balance:        balance,
		})
	}
	return monthlyPayment, payments
}
