package mortgage

import "github.com/spf13/cobra"

var MortgageCmd = &cobra.Command{
	Use:   "mortgage",
	Short: "Mortgage calculators",
}

var sep rune // fix this later

func init() {
	MortgageCmd.AddCommand(monthlyCmd)
	MortgageCmd.AddCommand(amortizeCmd)
}
