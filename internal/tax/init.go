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
var NewJerseyRegistry map[uint16]NjTaxSystem

func init() {
	// Get the directory of this source file
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)

	FicaRegistry = loadJson[map[uint16]FicaTaxSystem](filepath.Join(dir, "configs/fica.json"))
	UsFederalRegistry = loadJson[map[uint16]UsTaxSystem](filepath.Join(dir, "configs/us.json"))
	NewJerseyRegistry = loadJson[map[uint16]NjTaxSystem](filepath.Join(dir, "configs/nj.json"))
}

func loadJson[T any](filepath string) T {
	data, err := os.ReadFile(filepath)
	if err != nil {
		panic(fmt.Errorf("failed to read config file: %w", err))
	}

	var result T

	if err := json.Unmarshal(data, &result); err != nil {
		panic(fmt.Errorf("failed to unmarshal config: %w", err))
	}

	return result
}
