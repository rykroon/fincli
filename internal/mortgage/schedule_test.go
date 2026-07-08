package mortgage

import (
	"testing"

	"github.com/shopspring/decimal"
)

func dec(t *testing.T, s string) decimal.Decimal {
	t.Helper()
	d, err := decimal.NewFromString(s)
	if err != nil {
		t.Fatal(err)
	}
	return d
}

func TestCalculateMonthlyPayment(t *testing.T) {
	tests := []struct {
		name       string
		principal  string
		annualRate string
		years      uint16
		expected   string
	}{
		{"300k_7pct_30y", "300000", "0.07", 30, "1995.91"},
		{"200k_5pct_15y", "200000", "0.05", 15, "1581.59"},
		{"zero_rate", "120000", "0", 10, "1000.00"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			loan := NewLoan(dec(t, tc.principal), dec(t, tc.annualRate), tc.years)
			got := CalculateMonthlyPayment(loan.Principal, loan.MonthlyRate(), loan.NumPeriods())
			if !got.Round(2).Equal(dec(t, tc.expected)) {
				t.Errorf("got %s, want %s", got.Round(2), tc.expected)
			}
		})
	}
}

func TestCalculateSchedule(t *testing.T) {
	loan := NewLoan(dec(t, "300000"), dec(t, "0.07"), 30)
	sched := CalculateSchedule(loan, NewDefaultStrategy())

	if len(sched.Payments) != 360 {
		t.Errorf("expected 360 payments, got %d", len(sched.Payments))
	}

	finalBalance := sched.Payments[len(sched.Payments)-1].Balance
	if !finalBalance.Round(2).IsZero() {
		t.Errorf("expected final balance of 0, got %s", finalBalance.Round(2))
	}

	if !sched.TotalPrincipal.Round(2).Equal(loan.Principal) {
		t.Errorf(
			"total principal %s != loan principal %s",
			sched.TotalPrincipal.Round(2), loan.Principal,
		)
	}
}

func TestCalculateScheduleZeroRate(t *testing.T) {
	loan := NewLoan(dec(t, "120000"), dec(t, "0"), 10)
	sched := CalculateSchedule(loan, NewDefaultStrategy())

	if len(sched.Payments) != 120 {
		t.Errorf("expected 120 payments, got %d", len(sched.Payments))
	}
	if !sched.TotalInterest.IsZero() {
		t.Errorf("expected zero interest, got %s", sched.TotalInterest)
	}
}

func TestExtraPaymentsShortenSchedule(t *testing.T) {
	loan := NewLoan(dec(t, "300000"), dec(t, "0.07"), 30)
	zero := dec(t, "0")
	base := CalculateSchedule(loan, NewDefaultStrategy())
	extraMonthly := CalculateSchedule(loan, NewExtraPaymentStrategy(dec(t, "200"), zero))
	extraAnnual := CalculateSchedule(loan, NewExtraPaymentStrategy(zero, dec(t, "2400")))
	combined := CalculateSchedule(loan, NewExtraPaymentStrategy(dec(t, "200"), dec(t, "2400")))

	if len(extraMonthly.Payments) >= len(base.Payments) {
		t.Errorf(
			"extra monthly payments should shorten the schedule: %d >= %d",
			len(extraMonthly.Payments), len(base.Payments),
		)
	}
	if !extraMonthly.TotalInterest.LessThan(base.TotalInterest) {
		t.Error("extra monthly payments should reduce total interest")
	}

	if len(extraAnnual.Payments) >= len(base.Payments) {
		t.Errorf(
			"extra annual payments should shorten the schedule: %d >= %d",
			len(extraAnnual.Payments), len(base.Payments),
		)
	}

	if len(combined.Payments) >= len(extraMonthly.Payments) ||
		len(combined.Payments) >= len(extraAnnual.Payments) {
		t.Errorf(
			"combined extras should beat either alone: combined=%d monthly=%d annual=%d",
			len(combined.Payments), len(extraMonthly.Payments), len(extraAnnual.Payments),
		)
	}
	if !combined.TotalInterest.LessThan(extraMonthly.TotalInterest) {
		t.Error("combined extras should reduce interest below monthly-only")
	}
}

// A payment that doesn't cover interest must not loop forever; the schedule
// bails out instead of growing the balance indefinitely.
func TestCalculateScheduleTerminatesWhenPaymentTooSmall(t *testing.T) {
	loan := NewLoan(dec(t, "300000"), dec(t, "0.07"), 30)
	sched := CalculateSchedule(loan, NewExtraPaymentStrategy(dec(t, "-5000"), dec(t, "0")))
	if len(sched.Payments) > 1 {
		t.Errorf("expected schedule to bail out, got %d payments", len(sched.Payments))
	}
}
