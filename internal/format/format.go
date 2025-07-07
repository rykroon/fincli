package format

import (
	"strings"

	"github.com/shopspring/decimal"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func ToMoney(d decimal.Decimal) string {
	s := d.StringFixed(2)
	parts := strings.Split(s, ".")
	if len(parts) != 2 {
		return ""
	}
	lhs := parts[0]
	i := d.NumDigits()
	newString := ""
	for idx, chr := range lhs {
		if i%3 == 0 && idx != 0 {
			newString += ","
		}
		newString += string(chr)
		i -= 1
	}
	return "$" + newString + "." + parts[1]
}

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
