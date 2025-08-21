package cmd

import (
	"fmt"
	"os"

	"github.com/rykroon/fincli/cmd/mortgage"
	"github.com/rykroon/fincli/internal/flagx"
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
	rootCmd.AddCommand(mortgage.MortgageCmd)
	rootCmd.AddCommand(homeCmd)
	rootCmd.AddCommand(fireCmd)
	rootCmd.AddCommand(taxCmd)
	rootCmd.AddCommand(playgroundCmd)

	rootCmd.PersistentFlags().Var(flagx.NewRuneVal(&sep, []rune{',', '_'}), "sep", "thousands separator")
}
