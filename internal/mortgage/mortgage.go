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

type sumPaymentsResult struct {
	TotalPayments float64
	TotalInterest float64
}

func SumPayments(payments []payment) sumPaymentsResult {
	t := sumPaymentsResult{}
	for _, p := range payments {
		t.TotalInterest += p.InterestPaid
		t.TotalPayments += p.PrincipalPaid + p.InterestPaid
	}
	return t
}

func CalculatePayments(p float64, i float64, n int, extraPayment float64) (float64, []payment) {
	monthlyPayment := CalculateMonthlyPayment(p, i, n)
	payments := []payment{}
	balance := p
	for period := 1; math.Round(balance) > 0; period++ {
		interestPaid := balance * i
		principalPaid := monthlyPayment - interestPaid + extraPayment

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
