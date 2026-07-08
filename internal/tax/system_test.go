package tax

import (
	"testing"
)

func loadSystem(t *testing.T, year uint16, name string) TaxSystem {
	t.Helper()
	system, err := LoadTaxSystem(year, name)
	if err != nil {
		t.Fatal(err)
	}
	return system
}

func assertTaxes(t *testing.T, system TaxSystem, p TaxPayer, expected string) {
	t.Helper()
	result, err := system.CalculateTax(p)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Taxes.Round(2).Equal(dec(t, expected)) {
		t.Errorf("Taxes = %s, want %s", result.Taxes.Round(2), expected)
	}
}

func TestUsTaxSystem2025(t *testing.T) {
	system := loadSystem(t, 2025, "us")

	// single, 100k: taxable = 85,000
	// 11,925*0.10 + 36,550*0.12 + 36,525*0.22 = 13,614
	assertTaxes(t, system, NewTaxPayer(dec(t, "100000"), Single), "13614")

	// income below the standard deduction clamps to zero taxable income
	assertTaxes(t, system, NewTaxPayer(dec(t, "10000"), Single), "0")

	// 401k adjustment reduces AGI: 100k - 10k -> taxable 75,000
	// 11,925*0.10 + 36,550*0.12 + 26,525*0.22 = 11,414
	assertTaxes(t, system,
		NewTaxPayer(dec(t, "100000"), Single, Adjustment{Label: "401k", Amount: dec(t, "10000")}),
		"11414",
	)
}

func TestUsTaxSystemUnsupportedFilingStatus(t *testing.T) {
	system := loadSystem(t, 2025, "us")
	_, err := system.CalculateTax(NewTaxPayer(dec(t, "100000"), HeadOfHouse))
	if err == nil {
		t.Fatal("expected error for filing status with no config data")
	}
}

func TestNjTaxSystem2025(t *testing.T) {
	system := loadSystem(t, 2025, "nj")

	// married_joint, 75k:
	// 20,000*0.014 + 30,000*0.0175 + 20,000*0.0245 + 5,000*0.035 = 1,470
	assertTaxes(t, system, NewTaxPayer(dec(t, "75000"), MarriedJoint), "1470")

	// single, 50k: 20,000*0.014 + 15,000*0.0175 + 5,000*0.035 + 10,000*0.05525 = 1,270
	assertTaxes(t, system, NewTaxPayer(dec(t, "50000"), Single), "1270")
}

func TestFicaTaxSystem2025(t *testing.T) {
	system := loadSystem(t, 2025, "fica")

	// income above the SS wage base cap (176,100):
	// 176,100*0.062 + 200,000*0.0145 = 10,918.20 + 2,900 = 13,818.20
	assertTaxes(t, system, NewTaxPayer(dec(t, "200000"), Single), "13818.20")

	// below the cap: 100,000*(0.062+0.0145) = 7,650
	assertTaxes(t, system, NewTaxPayer(dec(t, "100000"), Single), "7650")
}

func TestLoadTaxSystemReturnsFreshInstances(t *testing.T) {
	a := loadSystem(t, 2024, "us")
	b := loadSystem(t, 2025, "us")
	if a == b {
		t.Fatal("LoadTaxSystem returned the same instance for different years")
	}

	// deductions differ between years, so shared state would corrupt results
	sdA := a.(*UsTaxSystem).FilingConfigs[Single].StandardDeduction
	sdB := b.(*UsTaxSystem).FilingConfigs[Single].StandardDeduction
	if sdA.Equal(sdB) {
		t.Errorf("2024 and 2025 standard deductions unexpectedly equal (%s)", sdA)
	}
}

func TestLoadTaxSystemErrors(t *testing.T) {
	if _, err := LoadTaxSystem(2025, "ny"); err == nil {
		t.Error("expected error for unsupported system")
	}
	if _, err := LoadTaxSystem(1999, "us"); err == nil {
		t.Error("expected error for missing year")
	}
}

func TestParseFilingStatus(t *testing.T) {
	for _, valid := range []string{"single", "married_joint", "married_separate", "head_of_household"} {
		if _, err := ParseFilingStatus(valid); err != nil {
			t.Errorf("ParseFilingStatus(%q) unexpected error: %v", valid, err)
		}
	}
	if _, err := ParseFilingStatus("bogus"); err == nil {
		t.Error("expected error for invalid filing status")
	}
}
