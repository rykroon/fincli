package mortgage

import (
	"github.com/spf13/cobra"
)

func NewMortgageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mortgage",
		Short: "Mortgage calculators",
	}

	cmd.AddCommand(NewMontlyCmd())
	cmd.AddCommand(NewAmortizeCmd())
	return cmd
}
