// movie_play.go — movie play <id>
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/movie-cli-v3/db"
)

var moviePlayCmd = &cobra.Command{
	Use:   "play [id]",
	Short: "Play a movie or TV show with the default player",
	Long:  `Opens a media file with the system's default video player.`,
	Args:  cobra.ExactArgs(1),
	Run:   runMoviePlay,
}

func runMoviePlay(cmd *cobra.Command, args []string) {
	database, err := db.Open()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Database error: %v\n", err)
		return
	}
	defer database.Close()

	id, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		fmt.Fprintln(os.Stderr, "❌ Invalid ID.")
		return
	}

	m, err := database.GetMediaByID(id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Media not found: %v\n", err)
		return
	}

	filePath := m.CurrentFilePath
	if _, statErr := os.Stat(filePath); statErr != nil {
		if os.IsNotExist(statErr) {
			fmt.Fprintf(os.Stderr, "❌ File not found: %s\n", filePath)
		} else {
			fmt.Fprintf(os.Stderr, "❌ Cannot access file %s: %v\n", filePath, statErr)
		}
		return
	}

	fmt.Printf("▶️  Playing: %s", m.CleanTitle)
	if m.Year > 0 {
		fmt.Printf(" (%d)", m.Year)
	}
	fmt.Println()

	var openCmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		openCmd = exec.Command("open", filePath)
	case "linux":
		openCmd = exec.Command("xdg-open", filePath)
	case "windows":
		openCmd = exec.Command("cmd", "/c", "start", "", filePath)
	default:
		fmt.Fprintf(os.Stderr, "❌ Unsupported OS: %s\n", runtime.GOOS)
		return
	}

	if err := openCmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Cannot open player: %v\n", err)
		return
	}
}
