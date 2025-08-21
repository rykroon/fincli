package fmtx

import "fmt"

func getFlags(state fmt.State) string {
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

func getWidth(state fmt.State) int {
	if w, ok := state.Width(); ok {
		return w
	}
	return -1
}

func getPrecision(state fmt.State) int {
	if p, ok := state.Precision(); ok {
		return p
	}
	return -1
}

func buildFormat(flags string, width, precision int, verb rune) string {
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
