package tax

import "github.com/shopspring/decimal"

type Stat struct {
	Name  string
	Value decimal.Decimal
}

type TaxResult struct {
	Name     string
	TaxesDue decimal.Decimal
	Stats    []Stat
}

func NewTaxResult(name string, total decimal.Decimal) TaxResult {
	return TaxResult{Name: name, TaxesDue: total, Stats: make([]Stat, 0)}
}

func (tr *TaxResult) AddStat(name string, value decimal.Decimal) {
	tr.Stats = append(tr.Stats, Stat{Name: name, Value: value})
}

type TaxSystem interface {
	CalculateTax(TaxPayer) TaxResult
}
