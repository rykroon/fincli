package realestate

import "github.com/spf13/cobra"

var RealEstateCmd = &cobra.Command{
	Use:   "real-estate",
	Short: "Real estate calculation.",
}

func init() {
	RealEstateCmd.AddCommand(housePurchaseCmd)
	RealEstateCmd.AddCommand(mortgagePayoffCmd)
}
