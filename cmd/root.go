package cmd

import (
	"fmt"
	"os"

	"github.com/rykroon/fincli/internal/fmtx"
	"github.com/spf13/cobra"
)

func Execute() {
	cmd := NewRootCmd()
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var prt fmtx.NumberPrinter

func NewRootCmd() *cobra.Command {
	sep := ","

	persistentPreRunE := func(cmd *cobra.Command, args []string) error {
		var sepRune rune
		switch sep {
		case ",":
			sepRune = ','
		case "_":
			sepRune = '_'
		case "none", "":
			sepRune = 0
		default:
			return fmt.Errorf("invalid value '%s' for sep, must be ',', '_', or 'none'", sep)
		}

		prt = fmtx.NewNumberPrinter(sepRune)
		return nil
	}

	cmd := &cobra.Command{
		Use:               "fin",
		Short:             "Finance CLI: Do Finance.",
		Long:              ``,
		PersistentPreRunE: persistentPreRunE,
	}

	cmd.AddCommand(NewMortgageCmd())
	cmd.AddCommand(NewHouseCmd())
	cmd.AddCommand(NewFireCmd())
	cmd.AddCommand(NewTaxCmd())

	cmd.PersistentFlags().StringVar(
		&sep, "sep", ",", "thousands separator: ',', '_', or 'none'",
	)
	return cmd
}
