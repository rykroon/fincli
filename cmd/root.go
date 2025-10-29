package cmd

import (
	"fmt"
	"os"

	"github.com/rykroon/fincli/cmd/mortgage"
	"github.com/rykroon/fincli/internal/flagx"
	"github.com/spf13/cobra"
)

var sep rune

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

	cmd.PersistentFlags().Var(flagx.NewRuneVal(&sep, []rune{',', '_'}), "sep", "thousands separator")
	return cmd
}
