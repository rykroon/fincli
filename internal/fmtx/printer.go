package fmtx

import (
	"fmt"

	"github.com/shopspring/decimal"
)

type NumberPrinter struct {
	sep rune
}

func NewNumberPrinter(sep rune) NumberPrinter {
	return NumberPrinter{sep}
}

func (p NumberPrinter) Printf(format string, a ...any) (int, error) {
	return fmt.Printf(format, p.transformArgs(a)...)
}

func (p NumberPrinter) Println(a ...any) (int, error) {
	return fmt.Println(a...)
}

// Wrap decimal.Decimal in DecFmt with separator.
func (p NumberPrinter) transformArgs(args []any) []any {
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
