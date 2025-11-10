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

func GetUsTaxSystem(year uint16) (*UsTaxSystem, error) {
	return getTaxSystem[UsTaxSystem](year, "us")
}

func GetFicaTaxSystem(year uint16) (*FicaTaxSystem, error) {
	return getTaxSystem[FicaTaxSystem](year, "fica")
}

func getTaxSystem[T any](year uint16, taxSystem string) (*T, error) {
	name := filepath.Join("configs", strconv.FormatUint(uint64(year), 10), taxSystem+".json")
	f, err := configs.Open(name)
	if err != nil {
		return nil, fmt.Errorf("failed to get tax system: %w", err)
	}
	defer f.Close()
	var result T
	err = json.NewDecoder(f).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("failed to get tax system: %w", err)
	}
	return &result, nil
}
