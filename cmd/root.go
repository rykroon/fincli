package cmd

import (
	"fmt"
	"os"

	"github.com/rykroon/fincli/cmd/mortgage"
	"github.com/rykroon/fincli/internal/flagx"
	"github.com/spf13/cobra"
)

func Execute() {
	cmd := NewRootCmd()
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fin",
		Short: "Finance CLI: Do Finance.",
		Long:  ``,
	}

	cmd.AddCommand(mortgage.NewMortgageCmd())
	cmd.AddCommand(NewHomeCmd())
	cmd.AddCommand(NewFireCmd())
	cmd.AddCommand(NewTaxCmd())

	var sep rune
	cmd.PersistentFlags().Var(flagx.NewRuneVal(&sep, []rune{',', '_'}), "sep", "thousands separator")
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
