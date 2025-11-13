package tax

import "github.com/shopspring/decimal"

type StatType string

const (
	Percent  StatType = "percent"
	Currency StatType = "currency"
)

type Stat struct {
	Name  string
	Value decimal.Decimal
	Type  StatType
}

type TaxResult struct {
	Name  string
	Taxes decimal.Decimal
	Stats []Stat
}

func NewTaxResult(name string, total decimal.Decimal) TaxResult {
	return TaxResult{Name: name, Taxes: total, Stats: make([]Stat, 0)}
}

func (tr *TaxResult) AddStat(name string, value decimal.Decimal, statType StatType) {
	tr.Stats = append(tr.Stats, Stat{Name: name, Value: value, Type: statType})
}

type TaxSystem interface {
	CalculateTax(TaxPayer) TaxResult
}
