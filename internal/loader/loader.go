package loader

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/LudovicMARIE/go-spotify-stats/internal/model"
)

func LoadAllPlays(dataDir string) ([]model.Play, error) {
	files, err := os.ReadDir(dataDir)
	if err != nil {
		return nil, fmt.Errorf("error reading data directory: %w", err)
	}

	var jsonFiles []string
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".json") {
			jsonFiles = append(jsonFiles, f.Name())
		}
	}

	if len(jsonFiles) == 0 {
		return nil, fmt.Errorf("no .json files found in directory: %s", dataDir)
	}

	var allPlays []model.Play
	playsChan := make(chan []model.Play, len(jsonFiles)) // Buffered channel to hold plays from each file
	errChan := make(chan error, len(jsonFiles))          // Channel to hold errors from each file
	done := make(chan bool)                              // Channel to signal completion of each goroutine

	for _, fileName := range jsonFiles {
		filePath := filepath.Join(dataDir, fileName)

		go func(fp string, fn string) {
			defer func() { done <- true }()

			plays, err := LoadTargetsFromFile(fp)
			if err != nil {
				errChan <- fmt.Errorf("error in file %s: %w", fn, err)
				return
			}
			playsChan <- plays
			errChan <- nil
		}(filePath, fileName)
	}

	go func() {
		for i := 0; i < len(jsonFiles); i++ {
			<-done
		}
		close(playsChan)
		close(errChan)
	}()

	var loadErrors []string
	for err := range errChan {
		if err != nil {
			loadErrors = append(loadErrors, err.Error())
		}
	}

	if len(loadErrors) > 0 {
		return nil, fmt.Errorf("errors during loading:\n%s", strings.Join(loadErrors, "\n"))
	}

	for plays := range playsChan {
		allPlays = append(allPlays, plays...)
	}

	return allPlays, nil
}
