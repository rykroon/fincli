package flagx

import (
	"fmt"
	"slices"
)

type RuneVal struct {
	ptr     *rune
	allowed []rune
}

func NewRuneVal(r *rune, allowed []rune) *RuneVal {
	return &RuneVal{r, allowed}
}

func (rv RuneVal) GetRune() rune {
	if rv.ptr == nil {
		return 0
	}
	return *rv.ptr
}

func (rv RuneVal) String() string {
	return string(rv.GetRune())
}

func (rv *RuneVal) Set(s string) error {
	if len(s) != 1 {
		return fmt.Errorf("must be a single character")
	}
	r := rune(s[0])

	if slices.Contains(rv.allowed, r) {
		*rv.ptr = r
		return nil
	}

	return fmt.Errorf("invalid separator: must be one of %q", string(rv.allowed))
}

func (rv RuneVal) Type() string {
	return "rune"
}
