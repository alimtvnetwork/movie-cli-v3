// movie_tag.go — movie tag [add|remove|list]
package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/movie-cli-v3/db"
	"github.com/spf13/cobra"
)

// movie tag — parent command
var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Manage tags on media items",
	Long:  "Add, remove, or list user-defined tags on movies and TV shows.",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

// movie tag add <id> <tag>
var tagAddCmd = &cobra.Command{
	Use:   "add <id> <tag>",
	Short: "Add a tag to a media item",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		d, err := db.Open()
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Database error: %v\n", err)
			return
		}
		defer d.Close()

		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Invalid ID: %s\n", args[0])
			return
		}
		tag := strings.TrimSpace(args[1])
		if tag == "" {
			fmt.Fprintln(os.Stderr, "❌ Tag cannot be empty")
			return
		}

		media, err := d.GetMediaByID(int64(id))
		if err != nil || media == nil {
			fmt.Fprintf(os.Stderr, "❌ Media not found with ID %d\n", id)
			return
		}

		err = d.AddTag(id, tag)
		if err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint") {
				fmt.Fprintf(os.Stderr, "❌ Tag \"%s\" already exists on \"%s\"\n", tag, media.Title)
			} else {
				fmt.Fprintf(os.Stderr, "❌ Error adding tag: %v\n", err)
			}
			return
		}

		year := ""
		if media.Year > 0 {
			year = fmt.Sprintf(" (%d)", media.Year)
		}
		fmt.Printf("✅ Tag \"%s\" added to \"%s%s\"\n", tag, media.Title, year)
	},
}

// movie tag remove <id> <tag>
var tagRemoveCmd = &cobra.Command{
	Use:   "remove <id> <tag>",
	Short: "Remove a tag from a media item",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		d, err := db.Open()
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Database error: %v\n", err)
			return
		}
		defer d.Close()

		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Invalid ID: %s\n", args[0])
			return
		}
		tag := strings.TrimSpace(args[1])

		media, err := d.GetMediaByID(int64(id))
		if err != nil || media == nil {
			fmt.Fprintf(os.Stderr, "❌ Media not found with ID %d\n", id)
			return
		}

		removed, err := d.RemoveTag(id, tag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Error removing tag: %v\n", err)
			return
		}
		if !removed {
			fmt.Fprintf(os.Stderr, "❌ Tag \"%s\" not found on \"%s\"\n", tag, media.Title)
			return
		}

		year := ""
		if media.Year > 0 {
			year = fmt.Sprintf(" (%d)", media.Year)
		}
		fmt.Printf("✅ Tag \"%s\" removed from \"%s%s\"\n", tag, media.Title, year)
	},
}

// movie tag list [id]
var tagListCmd = &cobra.Command{
	Use:   "list [id]",
	Short: "List tags (for a media item or all)",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		d, err := db.Open()
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Database error: %v\n", err)
			return
		}
		defer d.Close()

		if len(args) == 1 {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Fprintf(os.Stderr, "❌ Invalid ID: %s\n", args[0])
				return
			}

			media, err := d.GetMediaByID(int64(id))
			if err != nil || media == nil {
				fmt.Fprintf(os.Stderr, "❌ Media not found with ID %d\n", id)
				return
			}

			tags, err := d.GetTagsByMediaID(id)
			if err != nil {
				fmt.Fprintf(os.Stderr, "❌ Error reading tags: %v\n", err)
				return
			}

			year := ""
			if media.Year > 0 {
				year = fmt.Sprintf(" (%d)", media.Year)
			}

			if len(tags) == 0 {
				fmt.Printf("📭 No tags for \"%s%s\"\n", media.Title, year)
				return
			}

			fmt.Printf("🏷️  Tags for \"%s%s\":\n", media.Title, year)
			for _, t := range tags {
				fmt.Printf("  • %s\n", t)
			}
		} else {
			tagCounts, err := d.GetAllTagCounts()
			if err != nil {
				fmt.Fprintf(os.Stderr, "❌ Error reading tags: %v\n", err)
				return
			}

			if len(tagCounts) == 0 {
				fmt.Println("📭 No tags in library")
				return
			}

			fmt.Println("🏷️  All tags:")
			for _, tc := range tagCounts {
				fmt.Printf("  %s (%d)\n", tc.Tag, tc.Count)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(tagCmd)
	tagCmd.AddCommand(tagAddCmd)
	tagCmd.AddCommand(tagRemoveCmd)
	tagCmd.AddCommand(tagListCmd)
}
