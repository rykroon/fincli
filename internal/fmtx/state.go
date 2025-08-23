package fmtx

import (
	"fmt"
	"strings"
)

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
