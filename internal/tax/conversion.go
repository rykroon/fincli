package tax

import "github.com/shopspring/decimal"

type Number interface {
	int | float64 | decimal.Decimal
}

func numToDecimal[T Number](v T) decimal.Decimal {
	switch x := any(v).(type) {
	case int:
		return decimal.NewFromInt(int64(x))
	case float64:
		return decimal.NewFromFloat(x)
	case decimal.Decimal:
		return x
	default:
		panic("unexpected condition")
	}
}
