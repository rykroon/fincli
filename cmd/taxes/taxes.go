package taxes

import "github.com/spf13/cobra"

var TaxCmd = &cobra.Command{
	Use:   "taxes",
	Short: "Tax calculators",
}

func init() {
	// TaxCmd.AddCommand()
}
