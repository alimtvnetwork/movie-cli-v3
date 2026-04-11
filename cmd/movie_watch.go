// movie_watch.go — manage a personal watchlist (to-watch / watched).
package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/mahin/mahin-cli-v2/db"
)

var movieWatchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Manage your watchlist (to-watch / watched)",
	Long: `Track movies and TV shows you want to watch or have already watched.

Subcommands:
  movie watch add <ID>        Add a library item to your watchlist
  movie watch done <ID>       Mark as watched
  movie watch undo <ID>       Revert to to-watch
  movie watch rm <ID>         Remove from watchlist
  movie watch ls              List watchlist (default: to-watch)
  movie watch ls --all        List all entries
  movie watch ls --watched    List watched entries`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

var watchAddCmd = &cobra.Command{
	Use:   "add <media-ID>",
	Short: "Add a library item to your watchlist",
	Args:  cobra.ExactArgs(1),
	Run:   runWatchAdd,
}

var watchDoneCmd = &cobra.Command{
	Use:   "done <media-ID>",
	Short: "Mark a title as watched",
	Args:  cobra.ExactArgs(1),
	Run:   runWatchDone,
}

var watchUndoCmd = &cobra.Command{
	Use:   "undo <media-ID>",
	Short: "Revert a title back to to-watch",
	Args:  cobra.ExactArgs(1),
	Run:   runWatchUndo,
}

var watchRmCmd = &cobra.Command{
	Use:   "rm <media-ID>",
	Short: "Remove a title from your watchlist",
	Args:  cobra.ExactArgs(1),
	Run:   runWatchRm,
}

var (
	watchListAll     bool
	watchListWatched bool
)

var watchLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List your watchlist",
	Run:   runWatchLs,
}

func init() {
	watchLsCmd.Flags().BoolVar(&watchListAll, "all", false, "show all entries")
	watchLsCmd.Flags().BoolVar(&watchListWatched, "watched", false, "show watched entries only")

	movieWatchCmd.AddCommand(watchAddCmd, watchDoneCmd, watchUndoCmd, watchRmCmd, watchLsCmd)
}

func runWatchAdd(cmd *cobra.Command, args []string) {
	database, media := openAndGetMedia(args[0])
	if database == nil {
		return
	}
	defer database.Close()

	if err := database.AddToWatchlist(media.TmdbID, media.Title, media.Year, media.Type, media.ID); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
		return
	}
	fmt.Printf("📋 Added to watchlist: %s (%d)\n", media.Title, media.Year)
}

func runWatchDone(cmd *cobra.Command, args []string) {
	database, media := openAndGetMedia(args[0])
	if database == nil {
		return
	}
	defer database.Close()

	if err := database.MarkWatched(media.TmdbID); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
		return
	}
	fmt.Printf("✅ Marked as watched: %s (%d)\n", media.Title, media.Year)
}

func runWatchUndo(cmd *cobra.Command, args []string) {
	database, media := openAndGetMedia(args[0])
	if database == nil {
		return
	}
	defer database.Close()

	if err := database.MarkToWatch(media.TmdbID); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
		return
	}
	fmt.Printf("⏪ Reverted to to-watch: %s (%d)\n", media.Title, media.Year)
}

func runWatchRm(cmd *cobra.Command, args []string) {
	database, media := openAndGetMedia(args[0])
	if database == nil {
		return
	}
	defer database.Close()

	if err := database.RemoveFromWatchlist(media.TmdbID); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
		return
	}
	fmt.Printf("🗑️  Removed from watchlist: %s (%d)\n", media.Title, media.Year)
}

func runWatchLs(cmd *cobra.Command, args []string) {
	database, err := db.Open()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Database error: %v\n", err)
		return
	}
	defer database.Close()

	status := "to-watch"
	if watchListAll {
		status = ""
	} else if watchListWatched {
		status = "watched"
	}

	entries, err := database.ListWatchlist(status)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
		return
	}

	if len(entries) == 0 {
		label := "to-watch"
		if watchListAll {
			label = ""
		} else if watchListWatched {
			label = "watched"
		}
		if label != "" {
			fmt.Printf("📋 No %s entries in your watchlist.\n", label)
		} else {
			fmt.Println("📋 Your watchlist is empty.")
		}
		return
	}

	header := "📋 Watchlist — To Watch"
	if watchListAll {
		header = "📋 Watchlist — All"
	} else if watchListWatched {
		header = "📋 Watchlist — Watched"
	}
	fmt.Println(header)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	for i, e := range entries {
		icon := "🔲"
		if e.Status == "watched" {
			icon = "✅"
		}
		fmt.Printf("  %d. %s %s (%d) [%s]\n", i+1, icon, e.Title, e.Year, e.Type)
	}
}

// openAndGetMedia is a helper to open the DB and fetch media by ID arg.
func openAndGetMedia(idArg string) (*db.DB, *db.Media) {
	id, err := strconv.ParseInt(idArg, 10, 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Invalid ID: %s\n", idArg)
		return nil, nil
	}

	database, err := db.Open()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Database error: %v\n", err)
		return nil, nil
	}

	media, err := database.GetMediaByID(id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Media not found: %v\n", err)
		database.Close()
		return nil, nil
	}

	return database, media
}
