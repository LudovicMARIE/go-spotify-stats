package cmd

import (
	"fmt"
	"os"

	"github.com/LudovicMARIE/go-spotify-stats/internal/ingest"
	"github.com/LudovicMARIE/go-spotify-stats/internal/model"
	"github.com/LudovicMARIE/go-spotify-stats/internal/process"
	"github.com/spf13/cobra"
)

var DataDir string

var rootCmd = &cobra.Command{
	Use:   "go-spoti-stats",
	Short: "go spoti stats is a tool used to generate differents statistics based on your own spotify data",
	Long: `go spoti stats is a tool used to generate differents statistics based on your own spotify data.
	You can use it to analyze your listening habits, discover trends, and visualize your music preferences over time.
	All you need is to put your spotify data files (in JSON format) in a folder and point the tool to that folder.
	`,
	// Ensure the data directory exists before running any subcommand.
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if DataDir == "" {
			DataDir = "."
		}
		if _, err := os.Stat(DataDir); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "data directory %s does not exist\n", DataDir)
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Using data directory: %s\n", DataDir)

		files, err := os.ReadDir(DataDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading data directory: %v\n", err)
			os.Exit(1)
		}

		var allPlays []model.Play
		playsChan := make(chan []model.Play, len(files)) // Buffered channel to hold plays from each file
		errChan := make(chan error, len(files))          // Channel to hold errors from each file
		done := make(chan bool)                          // Channel to signal completion of each goroutine
		var numFiles int

		for _, f := range files {
			if f.IsDir() {
				continue
			}
			numFiles++
			filePath := DataDir + "/" + f.Name()
			go func(filePath string, fileName string) {
				defer func() { done <- true }()

				plays, err := ingest.LoadTargetsFromFile(filePath)
				if err != nil {
					errChan <- fmt.Errorf("error loading data from file %s: %v", fileName, err)
					return
				}
				playsChan <- plays
				errChan <- nil
			}(filePath, f.Name())
		}

		go func() {
			for i := 0; i < numFiles; i++ {
				<-done
			}
			close(playsChan)
			close(errChan)
		}()

		for err := range errChan {
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}

		for plays := range playsChan {
			allPlays = append(allPlays, plays...)
		}

		fmt.Printf("Loaded %d plays from %d files\n", len(allPlays), len(files))

		process.ProcessPlays(&allPlays)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func init() {
	// Add a persistent flag so every subcommand inherits it.
	rootCmd.PersistentFlags().StringVarP(&DataDir, "data-dir", "d", ".", "Directory to put the data files")
}
