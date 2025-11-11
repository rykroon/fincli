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

	run := func(cmd *cobra.Command, args []string) {
		taxPayer := tax.NewTaxPayer(
			income,
			tax.FilingStatus(filingStatus),
			tax.Adjustment{Label: "Adjustments", Amount: adjustments},
		)
		runTaxCmd(year, taxPayer)
	}

	cmd := &cobra.Command{
		Use:   "tax",
		Short: "Calculate Income Taxes",
		Run:   run,
	}

	flagx.DecimalVarP(cmd.Flags(), &income, "income", "i", decimal.Zero, "Your gross income")
	cmd.Flags().StringVarP(&filingStatus, "filing-status", "f", "single", "Your filing status")
	cmd.Flags().Uint16VarP(&year, "year", "y", 2025, "Tax year")
	flagx.DecimalVar(cmd.Flags(), &adjustments, "adjustments", decimal.Zero, "adjustments (ex: Retirement Contributions, Student Loan Interest)")
	cmd.MarkFlagRequired("income")
	return cmd
}

func runTaxCmd(year uint16, taxPayer tax.TaxPayer) error {
	prt.Printf("Gross Income: $%.2v\n", taxPayer.Income)
	prt.Println("")

	usTaxSystem, err := tax.LoadTaxSystem(year, "us")
	if err != nil {
		panic(err)
	}

	usTaxResult := usTaxSystem.CalculateTax(taxPayer)

	prt.Println("Federal Tax")
	for _, stat := range usTaxResult.Stats {
		prt.Printf("  %s: $%.2v\n", stat.Name, stat.Value)
	}

	prt.Printf("  Taxes Due: $%.2v\n", usTaxResult.TaxesDue)
	effectiveTaxRate := usTaxResult.TaxesDue.Div(taxPayer.Income)
	oneHundred := decimal.NewFromInt(100)
	prt.Printf("  Effective Tax Rate: %.2v%%\n", effectiveTaxRate.Mul(oneHundred))
	prt.Println("")

	// FICA Tax
	ficaTaxSystem, err := tax.LoadTaxSystem(year, "fica")
	if err != nil {
		panic("FICA tax system not found")
	}

	ficaTaxResult := ficaTaxSystem.CalculateTax(taxPayer)

	prt.Println("FICA Tax")
	for _, stat := range ficaTaxResult.Stats {
		prt.Printf("  %s: $%.2v\n", stat.Name, stat.Value)
	}
	prt.Println("")

	njTaxSystem, err := tax.LoadTaxSystem(year, "nj")
	if err != nil {
		panic("NJ tax system not found")
	}

	njTaxResult := njTaxSystem.CalculateTax(taxPayer)

	prt.Println("NJ Tax")
	for _, stat := range njTaxResult.Stats {
		prt.Printf("  %s: $%.2v\n", stat.Name, stat.Value)
	}

	prt.Printf("  Taxes Due: $%.2v", njTaxResult.TaxesDue)
	prt.Println("")

	return nil
}
