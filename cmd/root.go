package cmd

import (
	"fmt"
	"os"

	"github.com/rykroon/ry-cli/cmd/investing"
	"github.com/rykroon/ry-cli/cmd/realestate"
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
	rootCmd.AddCommand(realestate.RealEstateCmd)
	rootCmd.AddCommand(investing.InvestingCmd)
}
