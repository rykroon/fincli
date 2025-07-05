package finance

import (
	"fmt"

	"github.com/shopspring/decimal"
)

type Money struct {
	value decimal.Decimal
}

func (m *Money) Decimal() decimal.Decimal {
	return m.value
}

func (m1 Money) Add(m2 Money) Money {
	return Money{m1.Decimal().Add(m2.Decimal())}
}

func (m1 Money) Sub(m2 Money) Money {
	return Money{m1.Decimal().Sub(m2.Decimal())}
}

func (m *Money) String() string {
	return "$" + m.value.StringFixed(2)
}

func (m *Money) Set(s string) error {
	d, err := decimal.NewFromString(s)
	if err != nil {
		return fmt.Errorf("invalid decimal value: %w", err)
	}
	m.value = d
	return nil
}

func (m *Money) Type() string {
	return "decimal"
}
