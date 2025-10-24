package tax

type TaxRegistry[T TaxSystem] map[uint16]T
