package finance

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var FinanceCmd = &cobra.Command{
	Use:   "finance",
	Short: "Finance sub-command",
	Long:  ``,
	Run:   runFinanceCmd,
}

func Execute() {
	if err := FinanceCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runFinanceCmd(cmd *cobra.Command, args []string) {
	fmt.Println("Finance!")
}
