package fmtx

import (
	"fmt"

	"github.com/shopspring/decimal"
)

type DecimalPrinter struct {
	sep rune
}

func NewDecimalPrinter(sep rune) DecimalPrinter {
	return DecimalPrinter{sep}
}

func (p DecimalPrinter) Printf(format string, a ...any) (int, error) {
	return fmt.Printf(format, p.transformArgs(a)...)
}

func (p DecimalPrinter) Println(a ...any) (int, error) {
	return fmt.Println(a...)
}

// Wrap decimal.Decimal in DecFmt with separator.
func (p DecimalPrinter) transformArgs(args []any) []any {
	out := make([]any, 0, len(args))
	for _, arg := range args {
		switch v := arg.(type) {
		case decimal.Decimal:
			out = append(out, NewDecFmt(v, p.sep))
		default:
			out = append(out, v)
		}
	}
	return out
}
