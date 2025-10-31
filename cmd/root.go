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
	persistentPreRunE := func(cmd *cobra.Command, args []string) error {
		sep, err := flagx.GetRune(cmd.Flags(), "sep")
		if err != nil {
			return err
		}
		if sep != 0 && sep != ',' && sep != '_' {
			return fmt.Errorf("invalid value '%c' for sep, must be ',' or '_'", sep)
		}
		return nil
	}

	cmd := &cobra.Command{
		Use:               "fin",
		Short:             "Finance CLI: Do Finance.",
		Long:              ``,
		PersistentPreRunE: persistentPreRunE,
	}

	cmd.AddCommand(mortgage.NewMortgageCmd())
	cmd.AddCommand(NewHomeCmd())
	cmd.AddCommand(NewFireCmd())
	cmd.AddCommand(NewTaxCmd())
	flagx.Rune(cmd.PersistentFlags(), "sep", 0, "thousands separator")
	return cmd
}
