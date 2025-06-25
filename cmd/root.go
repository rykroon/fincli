package cmd

import (
	"fmt"
	"os"

	"github.com/rykroon/ry-cli/cmd/finance"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ry",
	Short: "Ry CLI is a command line tool for Ryan.",
	Long:  ``,
	Run:   runRootCmd,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runRootCmd(cmd *cobra.Command, args []string) {
	fmt.Println("root")
}

func init() {
	rootCmd.AddCommand(finance.FinanceCmd)
}
