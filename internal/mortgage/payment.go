package mortgage

import (
	"github.com/shopspring/decimal"
)

type Payment struct {
	Period    int
	Principal decimal.Decimal
	Interest  decimal.Decimal
	Balance   decimal.Decimal
}

func newPayment(
	period int,
	principal decimal.Decimal,
	interest decimal.Decimal,
	balance decimal.Decimal,
) Payment {
	return Payment{
		Period:    period,
		Principal: principal,
		Interest:  interest,
		Balance:   balance,
	}
}

func (p Payment) Total() decimal.Decimal {
	return p.Principal.Add(p.Interest)
}
