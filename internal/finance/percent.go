package finance

import (
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
)

type Percent struct {
	value decimal.Decimal
}

func (p *Percent) Decimal() decimal.Decimal {
	return p.value
}

func (p *Percent) ApplyTo(m Money) Money {
	return Money{m.Decimal().Mul(p.Decimal()).Round(2)}
}

func (p *Percent) String() string {
	return p.value.Mul(decimal.NewFromInt(100)).StringFixed(2) + "%"
}

func (p *Percent) Set(s string) error {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "%")
	d, err := decimal.NewFromString(s)
	if err != nil {
		return fmt.Errorf("invalid decimal value: %w", err)
	}
	p.value = d.Div(decimal.NewFromInt(100))
	return nil
}

func (p *Percent) Type() string {
	return "percent"
}
