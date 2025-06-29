package investing

import "github.com/spf13/cobra"

var InvestingCmd = &cobra.Command{
	Use:   "invest",
	Short: "Investment calculators.",
}

func init() {
	InvestingCmd.AddCommand(fireNumberCmd)
}
