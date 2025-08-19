package cli

import (
	"fmt"

	"github.com/shopspring/decimal"
)

/**
Idea, wrap decimal in Formatter type.
%p prints as percent.

*/

type Printer interface {
	Print(a ...any) (int, error)
	Printf(format string, a ...any) (int, error)
	Println(a ...any) (int, error)
}

type DecimalPrinter struct {
	Sep rune
}

func (p DecimalPrinter) tranformArgs(args []any) []any {
	out := make([]any, len(args))
	for _, arg := range args {
		switch v := arg.(type) {
		case decimal.Decimal:
			out = append(out, FormatDecimal(v, p.Sep)) // possibly wrap decimal in custom formatter type

		default:
			out = append(out, arg)
		}
	}
	return out
}

func (p DecimalPrinter) Print(a ...any) (int, error) {
	return fmt.Print(p.tranformArgs(a)...)
}

func (p DecimalPrinter) Printf(format string, a ...any) (int, error) {
	return fmt.Printf(format, p.tranformArgs(a)...)
}

func (p DecimalPrinter) Println(a ...any) (int, error) {
	return fmt.Println(p.tranformArgs(a)...)
}
