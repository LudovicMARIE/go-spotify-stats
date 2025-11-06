package ingest

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/LudovicMARIE/go-spotify-stats/internal/model"
)

func LoadTargetsFromFile(path string) ([]model.Play, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("couldn't read file %s : %w", path, err)
	}

	var targets []model.Play
	if err := json.Unmarshal(data, &targets); err != nil {
		return nil, fmt.Errorf("error while unserializing %s: %w", path, err)
	}
	return targets, nil

}

func SaveTargetsToFile(filePath string, targets []model.Play) error {
	data, err := json.MarshalIndent(targets, "", "  ")
	if err != nil {
		return fmt.Errorf("couldn't save file %s: %w", filePath, err)
	}
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("couldn't save file %s: %w", filePath, err)
	}
	return nil
}
