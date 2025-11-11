package tax

import (
	"embed"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strconv"
)

//go:embed configs
var configs embed.FS

var taxSystems map[string]TaxSystem = map[string]TaxSystem{
	"us":   &UsTaxSystem{},
	"fica": &FicaTaxSystem{},
	"nj":   &NjTaxSystem{},
}

func LoadTaxSystem(year uint16, name string) (TaxSystem, error) {
	system, ok := taxSystems[name]
	if !ok {
		return nil, fmt.Errorf("tax table for state not found")
	}
	filename := filepath.Join("configs", strconv.FormatUint(uint64(year), 10), name+".json")
	f, err := configs.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to load tax table: %w", err)
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(system)
	if err != nil {
		return nil, fmt.Errorf("failed to load tax table: %w", err)
	}
	return system, nil
}
