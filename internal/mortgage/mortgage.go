package mortgage

import (
	"math"
)

func CalculateMonthlyPayment(p float64, i float64, n int) float64 {
	return p * ((i * math.Pow(1+i, float64(n))) / (math.Pow(1+i, float64(n)) - 1))
}

type payment struct {
	Period        int
	PrincipalPaid float64
	InterestPaid  float64
	Balance       float64
}

func (p *payment) Payment() float64 {
	return p.PrincipalPaid + p.InterestPaid
}

type paymentScheduleStats struct {
	TotalPayments         float64
	TotalInterest         float64
	AverageMonthlyPayment float64
}

func GetPaymentScheduleStats(payments []payment) paymentScheduleStats {
	pss := paymentScheduleStats{}

	for _, p := range payments {
		pss.TotalInterest += p.InterestPaid
		pss.TotalPayments += p.PrincipalPaid + p.InterestPaid
	}
	pss.AverageMonthlyPayment = pss.TotalPayments / float64(len(payments))
	return pss
}

func CalculatePayments(p float64, i float64, n int, extraPaymentStratgey ExtraPaymentStrategy) (float64, []payment) {
	monthlyPayment := CalculateMonthlyPayment(p, i, n)
	payments := []payment{}
	balance := p
	for period := 1; math.Round(balance) > 0; period++ {
		interestPaid := balance * i
		principalPaid := monthlyPayment - interestPaid

		extraPayment := extraPaymentStratgey(period, principalPaid, interestPaid)
		principalPaid += extraPayment

		if principalPaid > balance {
			principalPaid = balance
		}

		balance -= principalPaid
		payments = append(payments, payment{
			Period:        period,
			InterestPaid:  interestPaid,
			PrincipalPaid: principalPaid,
			Balance:       balance,
		})
	}
	return monthlyPayment, payments
}
