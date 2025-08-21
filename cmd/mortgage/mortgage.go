package mortgage

import (
	"github.com/rykroon/fincli/internal/flagx"
	"github.com/spf13/cobra"
)

var MortgageCmd = &cobra.Command{
	Use:   "mortgage",
	Short: "Mortgage calculators",
}

func init() {
	MortgageCmd.AddCommand(monthlyCmd)
	MortgageCmd.AddCommand(amortizeCmd)
}

func getSep(cmd *cobra.Command) rune {
	flagPtr := cmd.Flags().Lookup("sep")
	if flagPtr == nil {
		return 0
	}
	runeVal, ok := flagPtr.Value.(*flagx.RuneVal)
	if !ok {
		return 0
	}

	return runeVal.GetRune()
}
