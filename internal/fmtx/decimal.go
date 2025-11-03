package fmtx

import (
	"fmt"
	"strconv"
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

func (fd FormattableDecimal) IsNegative() bool {
	return fd.value.IsNegative()
}

func (fd FormattableDecimal) IntPart() string {
	return addThousandsSep(uint64(fd.value.Abs().IntPart()), fd.sep)
}

func (f FormattableDecimal) FracPart(precision int) string {
	var s string
	if precision >= 0 {
		s = f.value.StringFixed(int32(precision))
	} else {
		s = f.value.String()
	}
	parts := strings.SplitN(s, ".", 2)
	if len(parts) > 1 {
		return "." + parts[1]
	}
	return ""
}

func (fd FormattableDecimal) Format(state fmt.State, verb rune) {
	switch verb {
	case 'v', 's':
		fmt.Fprint(state, FormatNumber(fd, state))

	case 'e', 'E', 'f', 'F', 'g', 'G':
		// handle decimal as float
		fmt.Fprintf(state, fmt.FormatString(state, verb), fd.value.InexactFloat64())

	default:
		fmt.Fprintf(state, fmt.FormatString(state, verb), fd.value)
	}
}

func addThousandsSep(num uint64, sep rune) string {
	s := strconv.FormatUint(num, 10)

	if sep == 0 {
		return s
	}

	var bld strings.Builder

	n := len(s)
	firstGroupLen := n % 3
	if firstGroupLen == 0 {
		firstGroupLen = 3
	}

	bld.WriteString(s[:firstGroupLen])
	for i := firstGroupLen; i < n; i += 3 {
		bld.WriteRune(sep)
		bld.WriteString(s[i : i+3])
	}
	return bld.String()
}
