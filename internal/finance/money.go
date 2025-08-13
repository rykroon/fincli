package finance

import (
	"github.com/shopspring/decimal"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type Money struct {
	decimal.Decimal
}

func (m Money) String() string {
	p := message.NewPrinter(language.English)
	f, _ := m.Decimal.Round(2).Float64()
	return p.Sprintf("$%.2f", f)
}

func (m *Money) Set(s string) error {
	d, err := decimal.NewFromString(s)
	if err != nil {
		return err
	}
	m.Decimal = d
	return nil
}

func (m Money) Type() string {
	return "money"
}
