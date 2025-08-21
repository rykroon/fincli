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
	w := getWidth(state)
	p := getPrecision(state)
	flags := getFlags(state)
	format := buildFormat(flags, w, p, verb)

	switch verb {
	case 'v', 's':
		format := buildFormat(flags, w, -1, 's')
		fmt.Fprintf(state, format, FormatDecimal(df.decimal, state.Flag('+'), int32(p), df.sep))

	case 'e', 'E', 'f', 'F', 'g', 'G':
		fmt.Fprintf(state, format, df.decimal.InexactFloat64())

	default:
		fmt.Fprintf(state, format, df.decimal)
	}
}

func FormatDecimal(d decimal.Decimal, alwaysPrintSign bool, precision int32, sep rune) string {
	if precision >= 0 {
		return formatNumberString(d.StringFixed(precision), alwaysPrintSign, sep)
	} else {
		return formatNumberString(d.String(), alwaysPrintSign, sep)
	}
}

func formatNumberString(s string, alwaysPrintSign bool, sep rune) string {
	parts := strings.SplitN(s, ".", 2)
	intPart := parts[0]
	sign := ""
	if strings.HasPrefix(intPart, "-") {
		sign = "-"
		intPart = intPart[1:]
	} else if alwaysPrintSign {
		sign = "+"
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
