// root.go — defines the root cobra command and wires all subcommands together.
// The only logic here is registering child commands and calling Execute().
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/movie-cli-v3/version"
)

var rootCmd = &cobra.Command{
	Use:   "movie",
	Short: "Movie CLI — manage your movie & TV show library",
	Long: fmt.Sprintf(`movie-cli %s — Movie & TV Show Library Manager

A cross-platform CLI tool for managing a personal movie and TV show
library. Scan local folders, clean filenames, fetch metadata from TMDb,
organize files, and track your collection — all from the terminal.

Quick Start:
  movie config set tmdb-key YOUR_API_KEY   Set your TMDb API key
  movie scan ~/Movies                       Scan a folder for videos
  movie ls                                  List your library
  movie search "Inception"                  Search TMDb
  movie info 1                              Show movie details

Management:
  movie move                                Move files interactively
  movie rename                              Batch-rename messy filenames
  movie undo                                Undo last move/rename
  movie play 1                              Play with default player

Discovery:
  movie suggest                             Get recommendations
  movie tag add 1 favorite                  Tag your movies
  movie stats                               Library statistics

System:
  movie version                             Show version info
  movie self-update                         Update via git pull

Documentation: https://github.com/alimtvnetwork/movie-cli-v3`, version.Short()),
	Version: version.Short(),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("movie-cli %s\n\n", version.Short())
		_ = cmd.Help()
	},
}

func init() {
	// Keep the CLI surface focused on project commands only.
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.SetVersionTemplate(fmt.Sprintf("movie-cli %s\n", version.Full()))
	rootCmd.AddCommand(
		helloCmd,
		versionCmd,
		selfUpdateCmd,
		movieScanCmd,
		movieLsCmd,
		movieSearchCmd,
		movieSuggestCmd,
		movieMoveCmd,
		movieUndoCmd,
		movieInfoCmd,
		moviePlayCmd,
		movieStatsCmd,
		movieRenameCmd,
		movieConfigCmd,
		movieExportCmd,
		movieDuplicatesCmd,
		movieCleanupCmd,
		movieWatchCmd,
		movieHistoryCmd,
	)
}

// Execute is called by main.go. It is the single public entry point.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
