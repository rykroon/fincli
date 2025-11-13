package cmd

import (
	"github.com/rykroon/fincli/internal/flagx"
	"github.com/rykroon/fincli/internal/tax"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
)

func NewTaxCmd() *cobra.Command {
	var income decimal.Decimal
	var filingStatus string
	var year uint16
	var adjustments decimal.Decimal
	var state string

	runE := func(cmd *cobra.Command, args []string) error {
		taxPayer := tax.NewTaxPayer(
			income,
			tax.FilingStatus(filingStatus),
			tax.Adjustment{Label: "Adjustments", Amount: adjustments},
		)

		systemNames := []string{"us", "fica"}
		if state != "" {
			systemNames = append(systemNames, state)
		}

		taxSystems := make([]tax.TaxSystem, 0, len(systemNames))

		for _, name := range systemNames {
			system, err := tax.LoadTaxSystem(year, name)
			if err != nil {
				return err
			}
			taxSystems = append(taxSystems, system)
		}

		runTaxCmd(taxPayer, taxSystems)
		return nil
	}

	cmd := &cobra.Command{
		Use:   "tax",
		Short: "Calculate Income Taxes",
		RunE:  runE,
	}

	flagx.DecimalVarP(cmd.Flags(), &income, "income", "i", decimal.Zero, "Your gross income")
	cmd.Flags().StringVarP(&filingStatus, "filing-status", "f", "single", "Your filing status")
	cmd.Flags().Uint16VarP(&year, "year", "y", 2025, "Tax year")
	cmd.Flags().StringVar(&state, "state", "", "State income tax")
	flagx.DecimalVar(cmd.Flags(), &adjustments, "adjustments", decimal.Zero, "adjustments (ex: Retirement Contributions, Student Loan Interest)")
	cmd.MarkFlagRequired("income")
	return cmd
}

func runTaxCmd(taxPayer tax.TaxPayer, systems []tax.TaxSystem) {
	prt.Printf("Gross Income: $%.2v\n", taxPayer.Income)
	prt.Printf("Filing Status: %s\n", taxPayer.FilingStatus)
	prt.Println("")

	totalTaxes := decimal.Zero
	oneHundred := decimal.NewFromInt(100)

	for _, system := range systems {
		result := system.CalculateTax(taxPayer)
		totalTaxes = totalTaxes.Add(result.Taxes)

		prt.Println(result.Name)
		for _, stat := range result.Stats {
			switch stat.Type {
			case "currency":
				prt.Printf("  %-22s: $%.2v\n", stat.Name, stat.Value)
			case "percent":
				prt.Printf("  %-22s: %.2v%%\n", stat.Name, stat.Value.Mul(oneHundred))
			}
		}
		prt.Println("")
	}

	prt.Printf("Total Taxes: $%.2v\n", totalTaxes)
}
