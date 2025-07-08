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
	digitStr := parts[0]
	i := d.NumDigits()
	moneyString := ""
	for idx, chr := range digitStr {
		if i%3 == 0 && idx != 0 {
			moneyString += ","
		}
		moneyString += string(chr)
		i -= 1
	}
	return "$" + moneyString + "." + parts[1]
}

func FormatMoney(d decimal.Decimal) string {
	p := message.NewPrinter(language.English)
	f, _ := d.Round(2).Float64()
	return p.Sprintf("$%.2f", f)
}
