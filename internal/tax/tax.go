package tax

type TaxSystem[T any] interface {
	CalculateTax(TaxPayer) T
}
