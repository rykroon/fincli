package tax

import (
	"embed"
	"encoding/json"
	"fmt"
	"path"
	"strconv"
)

//go:embed configs
var configs embed.FS

var taxSystemFactories = map[string]func() TaxSystem{
	"us":   func() TaxSystem { return &UsTaxSystem{} },
	"fica": func() TaxSystem { return &FicaTaxSystem{} },
	"nj":   func() TaxSystem { return &NjTaxSystem{} },
}

// stateSystems are the tax system names that may be passed as a state.
var stateSystems = map[string]bool{
	"nj": true,
}

func IsStateSystem(name string) bool {
	return stateSystems[name]
}

func LoadTaxSystem(year uint16, name string) (TaxSystem, error) {
	newSystem, ok := taxSystemFactories[name]
	if !ok {
		return nil, fmt.Errorf("tax system '%s' not supported", name)
	}
	system := newSystem()

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

	filename := path.Join("configs", yearStr, name+".json")

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
