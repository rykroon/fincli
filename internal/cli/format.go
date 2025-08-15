package cli

import (
	"strings"

	"github.com/shopspring/decimal"
)

func FormatPercent(d decimal.Decimal, places int32) string {
	return d.Mul(decimal.NewFromInt(100)).StringFixed(places) + "%"
}

func FormatMoney(d decimal.Decimal, sep rune) string {
	return "$" + FormatDecimal(d, sep)
}

func FormatDecimal(d decimal.Decimal, sep rune) string {
	s := d.StringFixed(2)
	if sep == 0 {
		return s
	}

	// Split into whole and fractional parts
	parts := strings.SplitN(s, ".", 2)
	intPart := parts[0]
	fracPart := ""
	if len(parts) > 1 {
		fracPart = parts[1]
	}

	// Add separators to the integer part
	var out strings.Builder
	n := len(intPart)
	for i, digit := range intPart {
		out.WriteRune(digit)
		if (n-i-1)%3 == 0 && i != n-1 {
			out.WriteRune(sep)
		}
	}

	// Append fractional part if present
	if fracPart != "" {
		out.WriteRune('.')
		out.WriteString(fracPart)
	}

	return out.String()
}
