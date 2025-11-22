package cmd

import (
	"fmt"
	"os"

	"github.com/LudovicMARIE/go-spotify-stats/internal/loader"
	"github.com/LudovicMARIE/go-spotify-stats/internal/model"
	"github.com/LudovicMARIE/go-spotify-stats/internal/ui"
	"github.com/spf13/cobra"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "launch go spoti stats as a terminal ui",
	Long:  `launch go spoti stats as a terminal ui`,
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
		fmt.Println("Loading data from:", DataDir)
		allPlays, err := loader.LoadAllPlays(DataDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading data: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully loaded %d plays.\n", len(allPlays))
		ExecuteTUI(&allPlays)
	},
}

func ExecuteTUI(allPlays *[]model.Play) {
	fmt.Println("TUI launched")
	// call the tview-based UI runner
	if err := ui.RunTUI(allPlays); err != nil {
		fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(tuiCmd)

	tuiCmd.PersistentFlags().StringVarP(&DataDir, "data-dir", "d", ".", "Directory to put the data files")

}
