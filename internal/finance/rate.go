package finance

import (
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
)

type RateFrequency string

const (
	Annual   RateFrequency = "annual"
	Monthly  RateFrequency = "monthly"
	Biweekly RateFrequency = "biweekly"
	Weekly   RateFrequency = "weekly"
)

type Rate struct {
	value     decimal.Decimal
	frequency RateFrequency
}

func NewFromFloat(f float64) Rate {
	return Rate{decimal.NewFromFloat(f), Annual}
}

func (r Rate) Decimal() decimal.Decimal {
	return r.value
}

func (r Rate) Annual() Rate {
	switch r.frequency {
	case Annual:
		return r
	case Monthly:
		return Rate{r.Decimal().Mul(decimal.NewFromInt(12)), Annual}
	case Biweekly:
		return Rate{r.Decimal().Mul(decimal.NewFromInt(26)), Annual}
	case Weekly:
		return Rate{r.Decimal().Mul(decimal.NewFromInt(52)), Annual}
	default:
		return r // probably not correct
	}
}

func (r Rate) Monthly() Rate {
	return Rate{r.Annual().Decimal().Div(decimal.NewFromInt(12)), Monthly}
}

func (r Rate) Biweekly() Rate {
	return Rate{r.Annual().Decimal().Div(decimal.NewFromInt(26)), Biweekly}
}

func (r Rate) ApplyTo(m Money) Money {
	return Money{m.Decimal().Mul(r.Annual().Decimal()).Round(2)}
}

func (r *Rate) String() string {
	return r.Annual().Decimal().Mul(decimal.NewFromInt(100)).StringFixed(2) + "%"
}

func (r *Rate) Set(s string) error {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "%")
	d, err := decimal.NewFromString(s)
	if err != nil {
		return fmt.Errorf("invalid decimal value: %w", err)
	}
	r.value = d.Div(decimal.NewFromInt(100))
	return nil
}

func (r *Rate) Type() string {
	return "rate"
}
