package amortization

import (
	"math"

	"github.com/spf13/cobra"
)

var AmortizationCmd = &cobra.Command{
	Use:   "amortization",
	Short: "Amortization Calculation.",
	Long:  ``,
}

type amortizationParams struct {
	amount      float64
	rate        float64
	years       int
	price       float64
	downPercent float64
}

func (ap *amortizationParams) Principal() float64 {
	if ap.amount > 0 {
		return ap.amount
	} else {
		return ap.price - ap.DownPayment()
	}
}

func (ap *amortizationParams) DownPayment() float64 {
	if ap.downPercent > 0 {
		return ap.price * ap.downPercent / 100
	}
	return 0
}

func (ap *amortizationParams) NumPeriods() int {
	return ap.years * 12
}

func (ap *amortizationParams) MonthlyRate() float64 {
	return ap.rate / 12 / 100
}

func init() {
	AmortizationCmd.AddCommand(paymentCmd)
	AmortizationCmd.AddCommand(scheduleCmd)
}

func addFlags(cmd *cobra.Command, ap *amortizationParams) {
	cmd.Flags().Float64VarP(&ap.amount, "amount", "a", 0, "The loan amount borrowed.")
	cmd.Flags().Float64VarP(&ap.rate, "rate", "r", 0, "Annual interest rate.")
	cmd.Flags().IntVarP(&ap.years, "years", "y", 30, "Loan term in years.")

	// price and down percent can be used instead of amount.
	cmd.Flags().Float64VarP(&ap.price, "price", "p", 0, "Home purchase price.")
	cmd.Flags().Float64VarP(&ap.downPercent, "down-percent", "d", 0, "Down payment as a percent.")

	cmd.MarkFlagsOneRequired("amount", "price")
	cmd.MarkFlagsRequiredTogether("price", "down-percent")
	cmd.MarkFlagsMutuallyExclusive("amount", "price")
	cmd.MarkFlagsMutuallyExclusive("amount", "down-percent")
	cmd.MarkFlagRequired("rate")
}

type Payment struct {
	Period        int
	PrincipalPaid float64
	InterestPaid  float64
	Balance       float64
}

func getAmortizedPayment(p float64, r float64, n int) float64 {
	return p * ((r * math.Pow(1+r, float64(n))) / (math.Pow(1+r, float64(n)) - 1))
}

func getAmortizationSchedule(p float64, r float64, n int) (float64, []Payment) {
	// add extra payment option.
	monthlyPayment := getAmortizedPayment(p, r, n)
	payments := []Payment{}
	balance := p
	for i := range n {
		interestPaid := balance * r
		principalPaid := monthlyPayment - interestPaid
		balance -= principalPaid
		payments = append(payments, Payment{
			Period:        i + 1,
			InterestPaid:  interestPaid,
			PrincipalPaid: principalPaid,
			Balance:       balance,
		})
		if balance <= 0 {
			break
		}
	}
	return monthlyPayment, payments
}
