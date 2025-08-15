package cli

import "fmt"

type runeValue struct {
	target  *rune
	allowed []rune
}

func RuneValue(r *rune, allowed []rune) *runeValue {
	return &runeValue{r, allowed}
}

func (rv *runeValue) Set(s string) error {
	if len(s) != 1 {
		return fmt.Errorf("must be a single character")
	}
	r := rune(s[0])

	for _, a := range rv.allowed {
		if r == a {
			*rv.target = r
			return nil
		}
	}

	return fmt.Errorf("invalid separator: must be one of %q", string(rv.allowed))
}

func (rv *runeValue) String() string {
	if rv.target == nil {
		return ""
	}
	return string(*rv.target)
}

func (rv *runeValue) Type() string {
	return "rune"
}
