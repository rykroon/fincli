package tax

import (
	"encoding/json"
	"testing"

	"github.com/shopspring/decimal"
)

func newTestSchedule(t *testing.T) ProgressiveTax {
	t.Helper()
	return ProgressiveTax{
		Brackets: []Bracket{
			{Lower: dec(t, "0"), Upper: dec(t, "10000"), Rate: dec(t, "0.10")},
			{Lower: dec(t, "10000"), Upper: dec(t, "50000"), Rate: dec(t, "0.20")},
			{Lower: dec(t, "50000"), Upper: dec(t, "999999999999"), Rate: dec(t, "0.30")},
		},
	}
}

func dec(t *testing.T, s string) decimal.Decimal {
	t.Helper()
	d, err := decimal.NewFromString(s)
	if err != nil {
		t.Fatal(err)
	}
	return d
}

func TestProgressiveTaxCalculateTax(t *testing.T) {
	schedule := newTestSchedule(t)

	tests := []struct {
		name     string
		income   string
		expected string
	}{
		{"zero_income", "0", "0"},
		{"negative_income", "-5000", "0"},
		{"within_first_bracket", "5000", "500"},
		{"first_bracket_boundary", "10000", "1000"},
		{"second_bracket", "30000", "5000"},          // 1000 + 20000*0.20
		{"second_bracket_boundary", "50000", "9000"}, // 1000 + 8000
		{"top_bracket", "100000", "24000"},           // 1000 + 8000 + 50000*0.30
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := schedule.CalculateTax(dec(t, tc.income))
			if !got.Equal(dec(t, tc.expected)) {
				t.Errorf("CalculateTax(%s) = %s, want %s", tc.income, got, tc.expected)
			}
		})
	}
}

func TestGetMarginalBracket(t *testing.T) {
	schedule := newTestSchedule(t)

	tests := []struct {
		name         string
		income       string
		expectedRate string
	}{
		{"negative_income", "-5000", "0.10"},
		{"zero_income", "0", "0.10"},
		{"first_bracket", "5000", "0.10"},
		{"second_bracket", "30000", "0.20"},
		{"top_bracket", "1000000", "0.30"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := schedule.GetMarginalBracket(dec(t, tc.income))
			if !got.Rate.Equal(dec(t, tc.expectedRate)) {
				t.Errorf(
					"GetMarginalBracket(%s).Rate = %s, want %s",
					tc.income, got.Rate, tc.expectedRate,
				)
			}
		})
	}
}

func TestGetMarginalBracketEmpty(t *testing.T) {
	schedule := ProgressiveTax{}
	got := schedule.GetMarginalBracket(dec(t, "100"))
	if !got.Rate.IsZero() {
		t.Errorf("expected zero bracket, got rate %s", got.Rate)
	}
}

// TestConfigBracketsAreContiguous loads every embedded tax config and asserts
// that each progressive schedule's brackets start at zero and are contiguous
// (each bracket's upper equals the next bracket's lower).
func TestConfigBracketsAreContiguous(t *testing.T) {
	years, err := configs.ReadDir("configs")
	if err != nil {
		t.Fatal(err)
	}

	for _, yearDir := range years {
		if !yearDir.IsDir() {
			continue
		}
		year := yearDir.Name()

		for _, name := range []string{"us", "nj"} {
			data, err := configs.ReadFile("configs/" + year + "/" + name + ".json")
			if err != nil {
				t.Fatalf("%s/%s: %v", year, name, err)
			}

			schedules := map[string]ProgressiveTax{}
			switch name {
			case "us":
				var sys UsTaxSystem
				if err := json.Unmarshal(data, &sys); err != nil {
					t.Fatalf("%s/%s: %v", year, name, err)
				}
				for status, config := range sys.FilingConfigs {
					schedules[string(status)] = config.Schedule
					if !config.StandardDeduction.IsPositive() {
						t.Errorf(
							"%s/%s %s: standard deduction must be positive",
							year, name, status,
						)
					}
				}
			case "nj":
				var sys NjTaxSystem
				if err := json.Unmarshal(data, &sys); err != nil {
					t.Fatalf("%s/%s: %v", year, name, err)
				}
				for status, schedule := range sys.FilingConfigs {
					schedules[string(status)] = schedule
				}
			}

			for status, schedule := range schedules {
				validateBrackets(t, year+"/"+name+" "+status, schedule.Brackets)
			}
		}
	}
}

func validateBrackets(t *testing.T, label string, brackets []Bracket) {
	t.Helper()
	if len(brackets) == 0 {
		t.Errorf("%s: no brackets", label)
		return
	}
	if !brackets[0].Lower.IsZero() {
		t.Errorf("%s: first bracket must start at 0, got %s", label, brackets[0].Lower)
	}
	for i := 0; i < len(brackets)-1; i++ {
		if !brackets[i].Upper.Equal(brackets[i+1].Lower) {
			t.Errorf(
				"%s: bracket %d upper (%s) != bracket %d lower (%s)",
				label, i, brackets[i].Upper, i+1, brackets[i+1].Lower,
			)
		}
	}
	for i, b := range brackets {
		if b.Upper.LessThanOrEqual(b.Lower) {
			t.Errorf("%s: bracket %d upper (%s) <= lower (%s)", label, i, b.Upper, b.Lower)
		}
		if b.Rate.IsNegative() || b.Rate.GreaterThanOrEqual(decimal.NewFromInt(1)) {
			t.Errorf("%s: bracket %d has implausible rate %s", label, i, b.Rate)
		}
	}
}
