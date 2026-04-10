// changelog.go — implements the `movie changelog` command.
// Reads CHANGELOG.md from the binary's directory (or repo root) and prints
// the latest version block or the full changelog.
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var changelogLatest bool

var changelogCmd = &cobra.Command{
	Use:   "changelog",
	Short: "Show the project changelog",
	Long: `Display the project changelog from CHANGELOG.md.

By default prints the full changelog. Use --latest to show only
the most recent version block.

Examples:
  movie changelog            Show full changelog
  movie changelog --latest   Show only the latest version`,
	Run: func(cmd *cobra.Command, args []string) {
		path := findChangelog()
		if path == "" {
			fmt.Fprintln(os.Stderr, "❌ CHANGELOG.md not found")
			os.Exit(1)
		}

		lines, err := readChangelog(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Failed to read CHANGELOG.md: %v\n", err)
			os.Exit(1)
		}

		if changelogLatest {
			lines = extractLatestBlock(lines)
		}

		for _, l := range lines {
			fmt.Println(l)
		}
	},
}

func init() {
	changelogCmd.Flags().BoolVar(&changelogLatest, "latest", false, "Show only the latest version block")
	rootCmd.AddCommand(changelogCmd)
}

// findChangelog looks for CHANGELOG.md relative to the binary, then CWD.
func findChangelog() string {
	// Try next to the binary
	exe, err := os.Executable()
	if err == nil {
		p := filepath.Join(filepath.Dir(exe), "CHANGELOG.md")
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}

	// Try current working directory
	if _, err := os.Stat("CHANGELOG.md"); err == nil {
		return "CHANGELOG.md"
	}

	return ""
}

// readChangelog reads the file into a slice of lines.
func readChangelog(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var lines []string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	return lines, sc.Err()
}

// extractLatestBlock returns lines from the first ## heading to the next ## heading.
func extractLatestBlock(lines []string) []string {
	var block []string
	inBlock := false

	for _, l := range lines {
		if strings.HasPrefix(l, "## ") {
			if inBlock {
				break // reached the next version heading
			}
			inBlock = true
		}
		if inBlock {
			block = append(block, l)
		}
	}

	return block
}
