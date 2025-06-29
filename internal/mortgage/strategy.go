package mortgage

type ExtraPaymentStrategy func(int, float64, float64) float64

func NoExtraPaymentStrategy() ExtraPaymentStrategy {
	return func(period int, principalPaid float64, interestPaid float64) float64 {
		return 0
	}
}

func NewExtraMonthlyPaymentStrategy(payment float64) ExtraPaymentStrategy {
	return func(period int, principalPaid float64, interestPaid float64) float64 {
		return payment
	}
}

func NewExtraAnnualPaymentStrategy(payment float64, startOfYear bool) ExtraPaymentStrategy {
	return func(period int, principalPaid float64, interestPaid float64) float64 {
		if (startOfYear && period%12 == 1) || (!startOfYear && period%12 == 0) {
			return payment
		}
		return 0
	}
}

func NewExtraMonthlyPaymentPerYear() ExtraPaymentStrategy {
	return func(period int, principalPaid float64, interestPaid float64) float64 {
		return (principalPaid + interestPaid) / 12
	}
}

func NewDoublePaymentStrategy() ExtraPaymentStrategy {
	return func(period int, principalPaid float64, interestPaid float64) float64 {
		return principalPaid + interestPaid
	}
}
