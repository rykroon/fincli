package cmd

import (
	"fmt"
	"os"

	"github.com/rykroon/fincli/internal/cli"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "fin",
	Short: "Finance CLI: Do Finance.",
	Long:  ``,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var sep rune

func init() {
	rootCmd.AddCommand(fireNumberCmd)
	rootCmd.AddCommand(taxCmd)
	rootCmd.AddCommand(homePurchCmd)
	rootCmd.AddCommand(mortgageCmd)

	rootCmd.PersistentFlags().Var(cli.RuneValue(&sep, []rune{',', '_'}), "sep", "thousands separator")
}
