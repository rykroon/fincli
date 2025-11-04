package mortgage

import (
	"github.com/rykroon/fincli/internal/flagx"
	"github.com/rykroon/fincli/internal/fmtx"
	"github.com/spf13/cobra"
)

var prt fmtx.NumberPrinter

func NewMortgageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mortgage",
		Short: "Mortgage calculators",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			sep, _ := flagx.GetRune(cmd.Flags(), "sep")
			prt = fmtx.NewNumberPrinter(sep)
		},
	}

	cmd.AddCommand(NewMontlyCmd())
	cmd.AddCommand(NewAmortizeCmd())
	return cmd
}
