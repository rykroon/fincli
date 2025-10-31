package flagx

import (
	"fmt"

	"github.com/spf13/pflag"
)

type runeValue rune

func newRuneValue(val rune, p *rune) *runeValue {
	*p = val
	return (*runeValue)(p)
}

func (r *runeValue) Set(s string) error {
	if len(s) != 1 {
		return fmt.Errorf("must be a single character")
	}
	*r = runeValue(s[0])
	return nil
}

func (r runeValue) String() string {
	return string(r)
}

func (r runeValue) Type() string { return "rune" }

func RuneVarP(f *pflag.FlagSet, p *rune, name, shorthand string, value rune, usage string) {
	f.VarP(newRuneValue(value, p), name, shorthand, usage)
}

func RuneVar(f *pflag.FlagSet, p *rune, name string, value rune, usage string) {
	RuneVarP(f, p, name, "", value, usage)
}

func Rune(f *pflag.FlagSet, name string, value rune, usage string) *rune {
	p := new(rune)
	RuneVar(f, p, name, value, usage)
	return p
}

func RuneP(f *pflag.FlagSet, name, shorthand string, value rune, usage string) *rune {
	p := new(rune)
	RuneVarP(f, p, name, shorthand, value, usage)
	return p
}

func GetRune(f *pflag.FlagSet, name string) (rune, error) {
	flag := f.Lookup("sep")
	if flag == nil {
		return 0, fmt.Errorf("flag '%s' not found", name)
	}
	value, ok := flag.Value.(*runeValue)
	if !ok {
		return 0, fmt.Errorf("flag '%s' not a rune", name)
	}
	return rune(*value), nil
}
