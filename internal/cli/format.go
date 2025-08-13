package cli

import (
	"github.com/shopspring/decimal"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type MoneyFormatter interface {
	FormatMoney(decimal.Decimal) string
}

type CommaFormatter struct {
	printer *message.Printer
}

func NewCommaFormatter() MoneyFormatter {
	return CommaFormatter{printer: message.NewPrinter(language.English)}
}

func (fmt CommaFormatter) FormatMoney(d decimal.Decimal) string {
	f, _ := d.Round(2).Float64()
	return fmt.printer.Sprintf("$%.2f", f)
}

type SansCommaFormatter struct{}

func NewSansCommaFormatter() MoneyFormatter {
	return SansCommaFormatter{}
}

func (f SansCommaFormatter) FormatMoney(d decimal.Decimal) string {
	return "$" + d.StringFixed(2)
}

func FormatPercent(d decimal.Decimal, places int32) string {
	return d.Mul(decimal.NewFromInt(100)).StringFixed(places) + "%"
}
