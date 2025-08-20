package fmtx

import "fmt"

func GetFlags(state fmt.State) string {
	flags := ""
	if state.Flag('+') {
		flags += "+"
	}
	if state.Flag('-') {
		flags += "-"
	}
	if state.Flag('0') {
		flags += "0"
	}
	if state.Flag(' ') {
		flags += " "
	}
	if state.Flag('#') {
		flags += "#"
	}

	return flags
}

func GetWidth(state fmt.State) int {
	if w, ok := state.Width(); ok {
		return w
	}
	return -1
}

func GetPrecision(state fmt.State) int {
	if p, ok := state.Precision(); ok {
		return p
	}
	return -1
}

func BuildFormatFromState(state fmt.State, verb rune) string {
	flags := ""

	// Collect possible flags
	if state.Flag('+') {
		flags += "+"
	}
	if state.Flag('-') {
		flags += "-"
	}
	if state.Flag('0') {
		flags += "0"
	}
	if state.Flag(' ') {
		flags += " "
	}
	if state.Flag('#') {
		flags += "#"
	}

	width, ok := state.Width()
	if !ok {
		width = -1
	}

	precision, ok := state.Precision()
	if !ok {
		precision = -1
	}

	return BuildFormat(flags, width, precision, verb)

}

func BuildFormat(flags string, width, precision int, verb rune) string {
	format := "%"
	format += flags
	if width >= 0 {
		format += fmt.Sprintf("%d", width)
	}
	if precision >= 0 {
		format += fmt.Sprintf(".%d", precision)
	}
	format += string(verb)
	return format
}
