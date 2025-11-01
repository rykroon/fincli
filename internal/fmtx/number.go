package fmtx

import "fmt"

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
