package cmd

import (
	"fmt"
	"os"

	"github.com/LudovicMARIE/go-spotify-stats/internal/loader"
	"github.com/LudovicMARIE/go-spotify-stats/internal/root"
	"github.com/spf13/cobra"
)

var DataDir string

var rootCmd = &cobra.Command{
	Use:   "",
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
		fmt.Println("Loading data from:", DataDir)
		allPlays, err := loader.LoadAllPlays(DataDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading data: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully loaded %d plays.\n", len(allPlays))
		root.RootMain()
	},
}

func ExecuteRoot() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func init() {
	// Add a persistent flag so every subcommand inherits it.
	rootCmd.PersistentFlags().StringVarP(&DataDir, "data-dir", "d", ".", "Directory to put the data files")

}
