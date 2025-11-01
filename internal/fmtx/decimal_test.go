package fmtx

import (
	"fmt"
	"testing"

	"github.com/shopspring/decimal"
)

func TestFormatDecimal(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		input    string
		sep      rune
		expected string
	}{
		// no thousands separator
		{"no_sep", "%v", "12345.67", 0, "12345.67"},

		// commas for thousands sep
		{"comma_sep", "%v", "12345.67", ',', "12,345.67"},

		// underscores for thousands sep
		{"underscore_sep", "%v", "12345.67", '_', "12_345.67"},

		// printing negative decimals
		{"negative", "%v", "-12345.67", ',', "-12,345.67"},

		// always printing a positive sign
		{"always_plus", "%+v", "12345.67", ',', "+12,345.67"},

		// leaving a space for a positive sign
		{"space_for_plus", "% v", "12345.67", ',', " 12,345.67"},

		// padding
		{"padding", "%10v", "123.45", 0, "    123.45"},

		// zero padding
		{"zero_padding", "%010v", "123.45", 0, "0000123.45"},

		// zero padding with positive sign
		{"zero_padding_with_pos_sign", "%+010v", "123.45", 0, "+000123.45"},

		// left align padding
		{"left_align", "%-10v", "123.45", 0, "123.45    "},

		// left align padding with zero pad (zero's should be ignored)
		{"left_align_with_zero_pad", "%0-10v", "123.45", 0, "123.45    "},

		// various precision
		{"precision2", "%.2v", "123.4567", 0, "123.46"},
		{"precision0", "%.0v", "123.4567", 0, "123"},
		{"precision_high", "%.5v", "123.4", 0, "123.40000"},

		// width + precision combined
		{"width_prec", "%12.2v", "123.4", 0, "      123.40"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			d, err := decimal.NewFromString(tc.input)
			if err != nil {
				t.Fatal(err)
			}
			decfmt := NewFormattableDecimal(d, tc.sep)
			got := fmt.Sprintf(tc.format, decfmt)
			if got != tc.expected {
				t.Errorf("Got %q, wanted %q", got, tc.expected)
			}
		})
	}
}
