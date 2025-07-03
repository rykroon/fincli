package mortgage

import "math"

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
