package invest

import "github.com/spf13/cobra"

var InvestCmd = &cobra.Command{
	Use:   "invest",
	Short: "Investment calculators.",
}

func init() {
	InvestCmd.AddCommand(fireNumberCmd)
}
