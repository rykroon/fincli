package fmtx

import (
	"fmt"
	"strings"
)

type FormattableNumber interface {
	Sign() int
	IntPart() string
	FracPart(int) string
}

func FormatNumber(num FormattableNumber, state fmt.State) string {
	sign := ""
	if num.Sign() < 0 {
		sign = "-"
	} else {
		if state.Flag('+') { // always print + sign
			sign = "+"
		} else if state.Flag(' ') { // leave space for sign
			sign = " "
		}
	}

	numStr := num.IntPart()
	precision, ok := state.Precision()
	if !ok {
		precision = -1
	}

	fracPart := num.FracPart(precision)
	if fracPart != "" {
		numStr += fracPart
	}

	strLen := len(sign) + len(numStr)
	padding := buildPadding(state, strLen)

	if state.Flag('-') { // left align
		return sign + numStr + padding
	} else if state.Flag('0') { // zero pad
		return sign + padding + numStr
	} else {
		return padding + sign + numStr
	}
}

func buildPadding(state fmt.State, strLen int) string {
	w, ok := state.Width()
	if !ok || strLen > w {
		return ""
	}
	padLen := w - strLen
	padChar := " "
	if state.Flag('0') && !state.Flag('-') {
		padChar = "0"
	}
	return strings.Repeat(padChar, padLen)
}
