package fmtx

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/shopspring/decimal"
)

type FormattableDecimal struct {
	value decimal.Decimal
	sep   rune
}

func NewFormattableDecimal(d decimal.Decimal, sep rune) FormattableDecimal {
	return FormattableDecimal{d, sep}
}

func (fd FormattableDecimal) Format(state fmt.State, verb rune) {
	switch verb {
	case 'v', 's':
		fmt.Fprint(state, FormatDecimal(state, fd.value, fd.sep))

	case 'e', 'E', 'f', 'F', 'g', 'G':
		// handle decimal as float
		fmt.Fprintf(state, fmt.FormatString(state, verb), fd.value.InexactFloat64())

	default:
		fmt.Fprintf(state, fmt.FormatString(state, verb), fd.value)
	}
}

func FormatDecimal(state fmt.State, d decimal.Decimal, sep rune) string {
	// --- sign ---
	var sign string
	if d.IsNegative() {
		sign = "-"
	} else {
		sign = GetPositiveSign(state)
	}

	// number string (with thousands sep)
	numStr := AddThousandsSep(d.Abs().BigInt(), sep)

	// precision (fractional part)
	var decStr string
	if p, ok := state.Precision(); ok {
		decStr = d.StringFixed(int32(p))
	} else {
		decStr = d.String()
	}
	parts := strings.SplitN(decStr, ".", 2)
	if len(parts) > 1 {
		numStr += "." + parts[1]
	}

	// --- width/ padding ---
	numLength := len(sign) + len(numStr)
	padding := BuildPadding(state, numLength)

	if LeftAlign(state) {
		// sign + num + padding
		return sign + numStr + padding
	} else if ZeroPad(state) {
		// sign + zeros + num
		return sign + padding + numStr
	} else {
		// spaces + sign + number
		return padding + sign + numStr
	}
}

func AddThousandsSep(b *big.Int, sep rune) string {
	intStr := b.String()
	if sep == 0 {
		return intStr
	}

	var bld strings.Builder

	if intStr[0] == '-' {
		intStr = intStr[1:]
		bld.WriteRune('-')
	}

	n := len(intStr)
	firstGroupLen := n % 3
	if firstGroupLen == 0 {
		firstGroupLen = 3
	}

	bld.WriteString(intStr[:firstGroupLen])
	for i := firstGroupLen; i < n; i += 3 {
		bld.WriteRune(sep)
		bld.WriteString(intStr[i : i+3])
	}
	return bld.String()
}
