package investing

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var rebalanceCmd = &cobra.Command{
	Use:   "rebalance",
	Short: "Rebalance portfolio.",
	Long:  ``,
	RunE:  runRebalanceCmd,
}

type rebalanceParams struct {
	Totals    []float64
	RawSlices []string
}

type Slice struct {
	Name       string
	Allocation int64
}

var rbp rebalanceParams

func runRebalanceCmd(cmd *cobra.Command, args []string) error {
	rawSlices, _ := cmd.Flags().GetStringArray("slice")
	slices := []Slice{}
	totalAllocation := int64(0)
	for _, s := range rawSlices {
		parts := strings.Split(s, ":")
		if len(parts) != 2 {
			return fmt.Errorf("invalid format for slice. Expected format: 'NAME:ALLOCATION'")
		}
		allocation, err := strconv.ParseInt(parts[1], 10, 8)
		if err != nil {
			return fmt.Errorf("invalid format for slice. Expected format: 'NAME:ALLOCATION'")
		}
		slices = append(slices, Slice{parts[0], allocation})
		totalAllocation += allocation
	}

	if totalAllocation > 100 {
		return fmt.Errorf("total allocation cannot be greater than 100")
	}

	total := float64(0)
	for _, t := range rbp.Totals {
		total += t
	}

	prntr := message.NewPrinter(language.English)

	prntr.Printf("%-6s %-12s\n", "Slice", "Amount")
	fmt.Println(strings.Repeat("-", 21))
	for _, s := range slices {
		prntr.Printf("%-6s $%-12.f\n", s.Name, total*float64(s.Allocation)/100)
	}

	return nil
}

func init() {
	rebalanceCmd.Flags().StringArrayVarP(&rbp.RawSlices, "slice", "s", []string{}, "Format: Name:Allocation")
	rebalanceCmd.Flags().Float64SliceVarP(&rbp.Totals, "total", "t", nil, "Total portfolio amount.")
}
