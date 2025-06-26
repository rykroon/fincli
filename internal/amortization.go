package internal

import "math"

type Payment struct {
	Month         int
	PrincipalPaid float64
	InterestPaid  float64
	Balance       float64
}

func AmortizeLoan(p float64, r float64, n int) float64 {
	return p * ((r * math.Pow(1+r, float64(n))) / (math.Pow(1+r, float64(n)) - 1))
}

func AmortizationSchedule(p float64, r float64, n int) []Payment {
	// add extra payment option.
	monthlyPayment := AmortizeLoan(p, r, n)
	payments := make([]Payment, n+1)
	balance := p
	// if using extra payment, use while loop
	// while balance > 0
	for i := range n {
		interestPaid := balance * r
		principalPaid := monthlyPayment - interestPaid
		balance -= principalPaid
		payments[i] = Payment{
			Month:         i + 1,
			InterestPaid:  interestPaid,
			PrincipalPaid: principalPaid,
			Balance:       balance,
		}
	}
	return payments
}
