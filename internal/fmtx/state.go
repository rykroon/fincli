package fmtx

import (
	"fmt"
	"strings"
)

// could probably be simplified even more by
// replacing Sign() with IsNegative()
// padding can be figured out by just knowing the length of the sign
// and length of the number.

type FormattableNumber interface {
	Sign(fmt.State) rune
	Number(fmt.State) string
	Padding(fmt.State) string
}

func FormatNumber(num FormattableNumber, state fmt.State) string {
	sign := num.Sign(state)
	numstr := num.Number(state)
	padding := num.Padding(state)

	if LeftAlign(state) {
		return string(sign) + numstr + padding
	} else if ZeroPad(state) {
		return string(sign) + padding + numstr
	} else {
		return padding + string(sign) + numstr
	}
}

// func GetWidth(state fmt.State) int {
// 	if w, ok := state.Width(); ok {
// 		return w
// 	}
// 	return -1
// }

// func GetPrecision(state fmt.State) int {
// 	if p, ok := state.Precision(); ok {
// 		return p
// 	}
// 	return -1
//}

func AlwaysPrintSign(state fmt.State) bool {
	return state.Flag('+')
}

func LeftAlign(state fmt.State) bool {
	return state.Flag('-')
}

func AlternateFormat(state fmt.State) bool {
	return state.Flag('#')
}

func LeaveSpaceForSign(state fmt.State) bool {
	return state.Flag(' ')
}

func ZeroPad(state fmt.State) bool {
	return state.Flag('0')
}

func GetPositiveSign(state fmt.State) string {
	if AlwaysPrintSign(state) {
		return "+"
	} else if LeaveSpaceForSign(state) {
		return " "
	} else {
		return ""
	}
}

func PadChar(state fmt.State) rune {
	if ZeroPad(state) && !LeftAlign(state) {
		return '0'
	}
	return ' '
}

func BuildPadding(state fmt.State, strLen int) string {
	w, ok := state.Width()
	if !ok || strLen > w {
		return ""
	}
	padLen := w - strLen
	return strings.Repeat(string(PadChar(state)), padLen)
}
