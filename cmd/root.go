package cmd

import (
	"fmt"
	"os"

	"github.com/rykroon/ry-cli/cmd/amortization"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ry",
	Short: "Ry CLI is a command line tool for Ryan.",
	Long:  ``,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(amortization.AmortizationCmd)
	rootCmd.AddCommand(rebalanceCmd)
}
