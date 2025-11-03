package fmtx

import (
	"fmt"
)

type FormattableNumber interface {
	IsNegative() bool
	IntPart() string
	FracPart(int) string
}

func FormatNumber(num FormattableNumber, state fmt.State) string {
	var sign string
	if num.IsNegative() {
		sign = "-"
	} else {
		sign = GetPositiveSign(state)
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
	padding := BuildPadding(state, strLen)

	if LeftAlign(state) {
		return sign + numStr + padding
	} else if ZeroPad(state) {
		return sign + padding + numStr
	} else {
		return padding + string(sign) + numStr
	}
}
