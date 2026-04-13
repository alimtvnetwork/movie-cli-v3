// update.go — implements the `movie self-update` command.
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/movie-cli-v3/updater"
)

var selfUpdateCmd = &cobra.Command{
	Use:     "self-update",
	Aliases: []string{"update"},
	Short:   "Update movie-cli to the latest version",
	Long: `Updates movie-cli by syncing the local source repository from GitHub.

Automatically finds the repository by checking:
  1. The directory where the binary is installed
  2. The current working directory
  3. A movie-cli-v3/ folder next to the binary

If no local repo is found, it clones a fresh copy next to the binary and
reports bootstrap success. After syncing, rebuild with: pwsh run.ps1`,
	Run: func(cmd *cobra.Command, args []string) {
		result, err := updater.Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Self-update failed: %v\n", err)
			os.Exit(1)
		}

		if result.Bootstrapped {
			fmt.Printf("\n✨ Bootstrapped local source repo in %s\n", result.RepoPath)
			fmt.Printf("🔁 Commit: %s\n", result.AfterCommit)
			fmt.Println("\n💡 Rebuild with: pwsh run.ps1")
			return
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
		fmt.Println("\n💡 Rebuild with: pwsh run.ps1")
	},
}
