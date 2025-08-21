package cmd

import (
	"fmt"

	"github.com/rykroon/fincli/internal/fmtx"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
)

var playgroundCmd = &cobra.Command{
	Use:   "playground",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		opts := fmtx.FormatOptions{
			AlwaysPrintSign: true,
			ThousandsSep:    '_',
			Precision:       2,
			Width:           12,
			ZeroPad:         true,
			LeftAlign:       true,
		}

		d := decimal.NewFromFloat(1234)

		fmt.Println(fmtx.FormatDecimal(d, opts), "hello world")
	},
}
