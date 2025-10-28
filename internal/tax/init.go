package tax

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

var UsFederalRegistry map[uint16]UsTaxSystem
var FicaRegistry map[uint16]FicaTaxSystem
var NewJerseyRegistry map[uint16]NjTaxSystem

func init() {
	FicaRegistry = loadJson[map[uint16]FicaTaxSystem]("fica.json")
	UsFederalRegistry = loadJson[map[uint16]UsTaxSystem]("us.json")
	NewJerseyRegistry = loadJson[map[uint16]NjTaxSystem]("nj.json")
}

func loadJson[T any](jsonfile string) T {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	jsonFilePath := filepath.Join(dir, "configs", jsonfile)

	data, err := os.ReadFile(jsonFilePath)
	if err != nil {
		panic(fmt.Errorf("failed to read config file: %w", err))
	}

	var result T

	if err := json.Unmarshal(data, &result); err != nil {
		panic(fmt.Errorf("failed to unmarshal config: %w", err))
	}

	return result
}
