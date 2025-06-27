package amortization

import (
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var scheduleCmd = &cobra.Command{
	Use:   "schedule",
	Short: "Print amortization schedule.",
	Long:  ``,
	Run:   runScheduleCmd,
}

var amp amortizationParams

func runScheduleCmd(cmd *cobra.Command, args []string) {
	fmt := message.NewPrinter(language.English)
	monthlyPayment, payments := getAmortizationSchedule(amp.Principal(), amp.MonthlyRate(), amp.NumPeriods())

	// header
	fmt.Printf("%-6s %-6s %-12s %-12s %-12s %-12s\n", "Year", "Period", "Payment", "Principal", "Interest", "Balance")
	fmt.Println(strings.Repeat("-", 55))

	for _, payment := range payments {
		// add cumalive interest and equity
		year := (payment.Period-1)/12 + 1
		fmt.Printf("%-6d %-6d $%-11.2f $%-11.2f $%-11.2f $%-11.2f\n",
			year, payment.Period, monthlyPayment, payment.PrincipalPaid, payment.InterestPaid, payment.Balance)
	}
}

func init() {
	addFlags(scheduleCmd, &amp)
	// add 'extra' flag for adding cumulative Principal, cumulative interest, and equity
	// add 'format' flag to print as csv, json, or ascii.
}
