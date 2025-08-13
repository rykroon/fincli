package home

import (
	"github.com/spf13/cobra"
)

var HomeCmd = &cobra.Command{
	Use:   "home",
	Short: "Home calculations",
}

func init() {
	HomeCmd.AddCommand(purchaseCmd)
	HomeCmd.AddCommand(mortgageCmd)
}
