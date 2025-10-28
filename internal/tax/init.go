package tax

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

const maxIncome uint64 = 1_000_000_000_000_000

var UsFederalRegistry map[uint16]UsTaxSystem
var FicaRegistry map[uint16]FicaTaxSystem

func init() {
	// Get the directory of this source file
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)

	FicaRegistry = loadFicaRegistry(dir)
	UsFederalRegistry = loadUsRegistry(dir)

	fmt.Println(FicaRegistry)
}

func loadFicaRegistry(dir string) map[uint16]FicaTaxSystem {
	// Build the path to the file in the same directory
	configPath := filepath.Join(dir, "configs/fica.json")

	data, err := os.ReadFile(configPath)
	if err != nil {
		panic(fmt.Errorf("failed to read config file: %w", err))
	}

	var registry map[uint16]FicaTaxSystem

	if err := json.Unmarshal(data, &registry); err != nil {
		panic(fmt.Errorf("failed to unmarshal config: %w", err))
	}

	return registry
}

func loadUsRegistry(dir string) map[uint16]UsTaxSystem {
	// Build the path to the file in the same directory
	configPath := filepath.Join(dir, "configs/us.json")

	data, err := os.ReadFile(configPath)
	if err != nil {
		panic(fmt.Errorf("failed to read config file: %w", err))
	}

	var registry map[uint16]UsTaxSystem

	if err := json.Unmarshal(data, &registry); err != nil {
		panic(fmt.Errorf("failed to unmarshal config: %w", err))
	}

	return registry
}

func buildNewJerseyTaxSystem() NjTaxSystem {
	sys := NewNjTaxSystem()
	sys.AddFilingConfig(Single, buildNjSingle2025())
	sys.AddFilingConfig(MarriedSeparate, buildNjSingle2025())
	sys.AddFilingConfig(MarriedJoint, buildNjMarried2025())
	sys.AddFilingConfig(HeadOfHouse, buildNjMarried2025())
	return sys
}

func buildNjSingle2025() ProgressiveTax {
	return NewProgressiveTax(
		NewBracket(0, 20_000, .014),
		NewBracket(20_000, 35_000, .0175),
		NewBracket(35_000, 40_000, .035),
		NewBracket(40_000, 75_000, .05525),
		NewBracket(75_000, 500_000, .0637),
		NewBracket(500_000, 1_000_000, .0897),
		NewBracket(1_000_000, maxIncome, .1075),
	)
}

func buildNjMarried2025() ProgressiveTax {
	return NewProgressiveTax(
		NewBracket(0, 20_000, .014),
		NewBracket(20_000, 50_000, .0175),
		NewBracket(50_000, 70_000, .0245),
		NewBracket(70_000, 80_000, .035),
		NewBracket(80_000, 150_000, .05525),
		NewBracket(150_000, 500_000, .0637),
		NewBracket(500_000, 1_000_000, .0897),
		NewBracket(1_000_000, maxIncome, .1075),
	)
}

var NewJerseyRegistry = map[uint16]NjTaxSystem{
	2025: buildNewJerseyTaxSystem(),
}
