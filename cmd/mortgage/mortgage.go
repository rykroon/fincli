package mortgage

import (
	"github.com/rykroon/fincli/internal/flagx"
	"github.com/spf13/cobra"
)

func NewMortgageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mortgage",
		Short: "Mortgage calculators",
	}

	cmd.AddCommand(NewMontlyCmd())
	cmd.AddCommand(NewAmortizeCmd())
	return cmd
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
