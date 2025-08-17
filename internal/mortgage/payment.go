package mortgage

import (
	"github.com/shopspring/decimal"
)

func CalculateMonthlyPayment(p decimal.Decimal, i decimal.Decimal, n decimal.Decimal) decimal.Decimal {
	// P * ((i * (1+i)^n) / ((1+i)^n - 1))
	one := decimal.NewFromInt(1)
	return p.Mul((i.Mul(i.Add(one).Pow(n))).Div(i.Add(one).Pow(n).Sub(one)))
}

type Payment struct {
	period         int
	principal      decimal.Decimal
	extraPrincipal decimal.Decimal
	interest       decimal.Decimal
	balancePrior   decimal.Decimal
}

func newPayment(
	period int,
	principal decimal.Decimal,
	interest decimal.Decimal,
	balancePrior decimal.Decimal,
) Payment {
	if principal.GreaterThan(balancePrior) {
		principal = balancePrior
	}
	return Payment{
		period:       period,
		principal:    principal,
		interest:     interest,
		balancePrior: balancePrior,
	}
}

// getters/ setters
func (p Payment) Period() int                     { return p.period }
func (p Payment) Principal() decimal.Decimal      { return p.principal }
func (p Payment) ExtraPrincipal() decimal.Decimal { return p.extraPrincipal }
func (p *Payment) SetExtraPrincipal(extra decimal.Decimal) {
	p.extraPrincipal = extra
	if p.TotalPrincipal().GreaterThan(p.BalancePrior()) {
		// make sure that adding the extra payment doesn't
		// payoff more than the balance.
		p.extraPrincipal = p.BalancePrior().Sub(p.Principal())
	}
}
func (p Payment) Interest() decimal.Decimal     { return p.interest }
func (p Payment) BalancePrior() decimal.Decimal { return p.balancePrior }

func (p Payment) Balance() decimal.Decimal {
	return p.BalancePrior().Sub(p.TotalPrincipal())
}

func (p Payment) TotalPrincipal() decimal.Decimal {
	return p.Principal().Add(p.ExtraPrincipal())
}

func (p Payment) Total() decimal.Decimal {
	return p.TotalPrincipal().Add(p.Interest())
}
