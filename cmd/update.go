// update.go — implements the `movie self-update` command.
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/mahin/mahin-cli-v2/updater"
)

var selfUpdateCmd = &cobra.Command{
	Use:     "self-update",
	Aliases: []string{"update"},
	Short:   "Pull latest files from the cloned git repository",
	Long: `self-update pulls the latest files from the current cloned git repository.

It runs:
  1. git rev-parse --show-toplevel
  2. git status --porcelain (must be clean)
  3. git pull --ff-only`,
	Run: func(cmd *cobra.Command, args []string) {
		result, err := updater.Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Self-update failed: %v\n", err)
			os.Exit(1)
		}

		if result.AlreadyLatest {
			fmt.Printf("✔ Already up to date (%s)\n", result.AfterCommit)
			return
		}

		fmt.Printf("\n✨ Pulled latest changes in %s\n", result.RepoPath)
		fmt.Printf("🔁 Commit: %s → %s\n", result.PreviousVersion, result.UpdatedTo)
		if result.Output != "" {
			fmt.Println(result.Output)
		}
	},
}
