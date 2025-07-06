package format

import (
	"github.com/shopspring/decimal"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func FormatMoney(d decimal.Decimal) string {
	p := message.NewPrinter(language.English)
	f, _ := d.Round(2).Float64()
	return p.Sprintf("$%.2f", f)
}

func FormatPercent(d decimal.Decimal) string {
	s := d.Mul(decimal.NewFromInt(100)).StringFixed(0)
	return s + "%"
}

func FormatRate(d decimal.Decimal) string {
	s := d.Mul(decimal.NewFromInt(100)).StringFixed(2)
	return s + "%"
}
