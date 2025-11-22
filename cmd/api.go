package cmd

import (
	"fmt"
	"os"

	"github.com/LudovicMARIE/go-spotify-stats/internal/loader"
	"github.com/LudovicMARIE/go-spotify-stats/internal/model"
	"github.com/spf13/cobra"
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "launch go spoti stats as an api",
	Long:  `launch go spoti stats as an api`,
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
		ExecuteApi(&allPlays)
	},
}

func ExecuteApi(allPlays *[]model.Play) {
	fmt.Println("Api launched")
}

func init() {
	rootCmd.AddCommand(apiCmd)

	apiCmd.PersistentFlags().StringVarP(&DataDir, "data-dir", "d", ".", "Directory to put the data files")

}
