package mortgage

type ExtraPaymentStrategy func(int, float64, float64) float64

func NoExtraPayment() ExtraPaymentStrategy {
	return func(period int, principalPaid float64, interestPaid float64) float64 {
		return 0
	}
}

func ExtraMonthlyPayment(payment float64) ExtraPaymentStrategy {
	return func(period int, principalPaid float64, interestPaid float64) float64 {
		return payment
	}
}

func ExtraAnnualPayment(payment float64) ExtraPaymentStrategy {
	return func(period int, principalPaid float64, interestPaid float64) float64 {
		if period%12 == 0 {
			return payment
		}
		return 0
	}
}

func ExtraMonthlyAndAnnualPayment(monthlyPayment float64, annualPayment float64) ExtraPaymentStrategy {
	return func(period int, principalPaid float64, interestPaid float64) float64 {
		payment := monthlyPayment
		if period%12 == 0 {
			payment += annualPayment
		}
		return payment
	}
}

func PrincipalMatchInterest() ExtraPaymentStrategy {
	return func(period int, principalPaid float64, interestPaid float64) float64 {
		if principalPaid < interestPaid {
			return interestPaid - principalPaid
		}
		return 0
	}
}
