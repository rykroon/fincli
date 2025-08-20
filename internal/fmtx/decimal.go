package fmtx

import (
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
)

type DecFmt struct {
	decimal decimal.Decimal
	sep     rune
}

func NewDecFmt(d decimal.Decimal, sep rune) DecFmt {
	return DecFmt{d, sep}
}

func (df DecFmt) Format(state fmt.State, verb rune) {
	w := GetWidth(state)
	p := GetPrecision(state)
	flags := GetFlags(state)

	switch verb {
	case 'v', 's':
		format := BuildFormat(flags, w, -1, 's')
		fmt.Fprintf(state, format, FormatDecimal(df.decimal, int32(p), df.sep))
	case 'f', 'F':
		fmt.Fprintf(state, BuildFormat(flags, w, p, verb), df.decimal.InexactFloat64())
	default:
		fmt.Fprintf(state, BuildFormat(flags, w, p, verb), df.decimal)
	}

}

func FormatDecimal(d decimal.Decimal, precision int32, sep rune) string {
	var s string
	if precision >= 0 {
		s = d.StringFixed(precision)
	} else {
		s = d.String()
	}

	parts := strings.SplitN(s, ".", 2)
	intPart := parts[0]
	sign := ""
	if strings.HasPrefix(intPart, "-") {
		sign = "-"
		intPart = intPart[1:]
	}

	// Insert separators
	var out strings.Builder
	for i, r := range intPart {
		if i > 0 && (len(intPart)-i)%3 == 0 {
			out.WriteRune(sep)
		}
		out.WriteRune(r)
	}

	if len(parts) == 2 {
		return sign + out.String() + "." + parts[1]
	}
	return sign + out.String()
}
