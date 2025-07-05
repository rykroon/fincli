package mortgage

import "github.com/shopspring/decimal"

type ExtraPaymentStrategy func(int, decimal.Decimal, decimal.Decimal) decimal.Decimal

func NoExtraPayment() ExtraPaymentStrategy {
	return func(period int, principal decimal.Decimal, interest decimal.Decimal) decimal.Decimal {
		return decimal.Zero
	}
}

func ExtraMonthlyPayment(payment decimal.Decimal) ExtraPaymentStrategy {
	return func(period int, principal decimal.Decimal, interest decimal.Decimal) decimal.Decimal {
		return payment
	}
}

func ExtraAnnualPayment(payment decimal.Decimal) ExtraPaymentStrategy {
	return func(period int, principal decimal.Decimal, interest decimal.Decimal) decimal.Decimal {
		if period%12 == 0 {
			return payment
		}
		return decimal.Zero
	}
}

func ExtraMonthlyAndAnnualPayment(monthlyPayment decimal.Decimal, annualPayment decimal.Decimal) ExtraPaymentStrategy {
	return func(period int, principal decimal.Decimal, interest decimal.Decimal) decimal.Decimal {
		payment := monthlyPayment
		if period%12 == 0 {
			payment = payment.Add(annualPayment)
		}
		return payment
	}
}
