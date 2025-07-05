package finance

import (
	"fmt"

	"github.com/shopspring/decimal"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type Money struct {
	value decimal.Decimal
}

func NewMoneyFromInt(i int64) Money {
	return Money{decimal.NewFromInt(i)}
}

func (m Money) Decimal() decimal.Decimal {
	return m.value
}

func (m1 Money) Add(m2 Money) Money {
	return Money{m1.Decimal().Add(m2.Decimal())}
}

func (m1 Money) Sub(m2 Money) Money {
	return Money{m1.Decimal().Sub(m2.Decimal())}
}

func (m1 Money) Mul(m2 Money) Money {
	return Money{m1.Decimal().Mul(m2.Decimal())}
}

func (m1 *Money) Div(m2 Money) Money {
	return Money{m1.Decimal().Div(m2.Decimal())}
}

func (m Money) GreaterThan(m2 Money) bool {
	return m.Decimal().GreaterThan(m2.Decimal())
}

func (m Money) String() string {
	p := message.NewPrinter(language.English)
	f, _ := m.Decimal().Round(2).Float64()
	return p.Sprintf("$%.2f", f)
}

func (m *Money) Set(s string) error {
	d, err := decimal.NewFromString(s)
	if err != nil {
		return fmt.Errorf("invalid decimal: %w", err)
	}
	m.value = d
	return nil
}

func (m Money) Type() string {
	return "money"
}
