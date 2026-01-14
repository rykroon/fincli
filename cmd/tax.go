package cmd

import (
	"strings"

	"github.com/rykroon/fincli/internal/flagx"
	"github.com/rykroon/fincli/internal/tax"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
)

func NewTaxCmd() *cobra.Command {
	income := decimal.Zero
	filingStatus := ""
	year := uint16(0)
	f01k := decimal.Zero
	fica := false
	state := ""

	runE := func(cmd *cobra.Command, args []string) error {
		taxPayer := tax.NewTaxPayer(
			income,
			tax.FilingStatus(filingStatus),
			tax.Adjustment{Label: "401k Contribution", Amount: f01k},
		)

		systemNames := []string{"us"}
		if fica {
			systemNames = append(systemNames, "fica")
		}
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

	cmd.Flags().VarP(flagx.NewDecimalFlag(&income), "income", "i", "Your gross income")
	cmd.Flags().StringVarP(&filingStatus, "filing-status", "f", "single", "Your filing status")
	cmd.Flags().Uint16VarP(&year, "year", "y", 2025, "Tax year")
	cmd.Flags().BoolVar(&fica, "fica", false, "Include FICA tax")
	cmd.Flags().StringVar(&state, "state", "", "State income tax")
	// flagx.DecimalVar(cmd.Flags(), &adjustments, "adjustments", decimal.Zero, "adjustments (ex: Retirement Contributions, Student Loan Interest)")
	cmd.Flags().Var(flagx.NewDecimalFlag(&f01k), "401k", "401k Contributions")
	cmd.MarkFlagRequired("income")
	return cmd
}

func runTaxCmd(taxPayer tax.TaxPayer, systems []tax.TaxSystem) {
	prt.Printf("%-20s $%12.2v\n", "Gross Income:", taxPayer.Income)
	prt.Printf("%-20s %-12s\n", "Filing Status:", taxPayer.FilingStatus)
	prt.Println("")

	totalTaxes := decimal.Zero
	oneHundred := decimal.NewFromInt(100)

	for _, system := range systems {
		result := system.CalculateTax(taxPayer)
		totalTaxes = totalTaxes.Add(result.Taxes)

		prt.Println(result.Name)
		prt.Println(strings.Repeat("-", 20))
		for _, stat := range result.Stats {
			switch stat.Type {
			case "currency":
				prt.Printf("%-20s $%12.2v\n", stat.Name+":", stat.Value)
			case "percent":
				prt.Printf("%-21s %12.2v%%\n", stat.Name+":", stat.Value.Mul(oneHundred))
			}
		}
		prt.Println("")
	}

	prt.Println("Total")
	prt.Println(strings.Repeat("-", 20))
	prt.Printf("%-20s $%12.2v\n", "Taxes:", totalTaxes)
	prt.Printf("%-21s %12.2v%%\n", "Effective Tax Rate:", totalTaxes.Div(taxPayer.Income).Mul(oneHundred))
	prt.Printf("%-20s $%12.2v\n", "Disposable Income:", taxPayer.Income.Sub(totalTaxes))
}
