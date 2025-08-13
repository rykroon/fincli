package cmd

import (
	"fmt"
	"os"

	"github.com/rykroon/fincli/cmd/home"
	"github.com/rykroon/fincli/cmd/invest"
	"github.com/rykroon/fincli/cmd/tax"
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

func init() {
	rootCmd.AddCommand(home.HomeCmd)
	rootCmd.AddCommand(invest.InvestCmd)
	rootCmd.AddCommand(tax.TaxCmd)
}
