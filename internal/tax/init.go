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
		return nil, fmt.Errorf("tax system '%s' not supported", name)
	}

	dirEntries, err := configs.ReadDir("configs")
	if err != nil {
		return nil, err
	}

	found := false
	yearStr := strconv.FormatUint(uint64(year), 10)
	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() && dirEntry.Name() == yearStr {
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("no tax tables found for year '%s'", yearStr)
	}

	filename := filepath.Join("configs", yearStr, name+".json")

	file, err := configs.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to load tax table: %w", err)
	}
	defer file.Close()
	err = json.NewDecoder(file).Decode(system)
	if err != nil {
		return nil, fmt.Errorf("failed to decode tax table: %w", err)
	}
	return system, nil
}
