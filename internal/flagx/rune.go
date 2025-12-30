package flagx

import (
	"fmt"

	"github.com/spf13/pflag"
)

func NewRuneFlag(r *rune) *Flag[rune] {
	return &Flag[rune]{
		ptr:        r,
		type_:      "rune",
		stringFunc: func(v rune) string { return string(v) },
		setFunc: func(s string, p *rune) error {
			if len(s) != 1 {
				return fmt.Errorf("must be a single character")
			}
			*r = rune(s[0])
			return nil
		},
	}
}

func GetRune(f *pflag.FlagSet, name string) (rune, error) {
	flag := f.Lookup(name)
	if flag == nil {
		return 0, fmt.Errorf("flag '%s' not found", name)
	}
	value, ok := flag.Value.(*Flag[rune])
	if !ok {
		return 0, fmt.Errorf("flag '%s' not a rune", name)
	}
	return rune(*value.ptr), nil
}
