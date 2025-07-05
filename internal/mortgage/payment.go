package mortgage

import (
	"github.com/shopspring/decimal"
)

func CalculateMonthlyPayment(p decimal.Decimal, i decimal.Decimal, n decimal.Decimal) decimal.Decimal {
	one := decimal.NewFromInt(1)
	return p.Mul((i.Mul(i.Add(one).Pow(n))).Div(i.Add(one).Pow(n).Sub(one)))
	// return p * ((i * math.Pow(1+i, float64(n))) / (math.Pow(1+i, float64(n)) - 1))
}

type Payment struct {
	Period         int
	Principal      decimal.Decimal
	ExtraPrincipal decimal.Decimal
	Interest       decimal.Decimal
	Balance        decimal.Decimal
}

func (p *Payment) TotalPrincipal() decimal.Decimal {
	return p.Principal.Add(p.ExtraPrincipal)
}

func (p *Payment) Total() decimal.Decimal {
	return p.TotalPrincipal().Add(p.Interest)
}
