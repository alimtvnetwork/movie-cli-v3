// movie_history.go — movie history: unified view of all tracked operations.
//
// Shows moves, renames, scans, deletions, popouts, and rescans from both
// move_history and action_history tables.
//
// Flags:
//
//	--type <type>    Filter by type: move, scan, delete, popout, rescan, all (default: all)
//	--batch <id>     Show all actions in a specific batch
//	--limit <n>      Max records to show (default: 20)
//	--format <fmt>   Output format: default, json, table
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/movie-cli-v3/db"
	"github.com/alimtvnetwork/movie-cli-v3/errlog"
)

var (
	historyFormat  string
	historyType    string
	historyBatch   string
	historySince   string
	historyLimit   int
)

var movieHistoryCmd = &cobra.Command{
	Use:   "history",
	Short: "Show history of all tracked operations",
	Long: `Displays the history of all state-changing operations including
file moves, renames, scans, deletions, popouts, and metadata rescans.

Flags:
  --type <type>   Filter: move, scan, delete, popout, rescan, all (default: all)
  --batch <id>    Show all actions in a specific batch
  --since <date>  Show only records after this date (e.g. 2026-04-01)
  --limit <n>     Max records (default: 20)
  --format <fmt>  Output: default, json, table`,
	Run: runMovieHistory,
}

func init() {
	movieHistoryCmd.Flags().StringVar(&historyFormat, "format", "default", "output format: default, json, table")
	movieHistoryCmd.Flags().StringVar(&historyType, "type", "all", "filter: move, scan, delete, popout, rescan, all")
	movieHistoryCmd.Flags().StringVar(&historyBatch, "batch", "", "show actions for a specific batch ID")
	movieHistoryCmd.Flags().StringVar(&historySince, "since", "", "show records after this date (e.g. 2026-04-01)")
	movieHistoryCmd.Flags().IntVar(&historyLimit, "limit", 20, "max records to show")
}

// unifiedRecord merges move_history and action_history into one display item.
type unifiedRecord struct {
	Source    string `json:"source"`     // "move" or "action"
	ID        int64  `json:"id"`
	Type      string `json:"type"`       // move, rename, scan_add, scan_remove, delete, popout, restore, rescan_update
	Detail    string `json:"detail"`
	FromPath  string `json:"from_path,omitempty"`
	ToPath    string `json:"to_path,omitempty"`
	BatchID   string `json:"batch_id,omitempty"`
	Timestamp string `json:"timestamp"`
	Undone    bool   `json:"undone"`
}

func runMovieHistory(cmd *cobra.Command, args []string) {
	database, err := db.Open()
	if err != nil {
		errlog.Error("Database error: %v", err)
		return
	}
	defer database.Close()

	// --batch: show a specific batch
	if historyBatch != "" {
		showBatchHistory(database)
		return
	}

	records := collectUnifiedRecords(database)

	if len(records) == 0 {
		if historyFormat == "json" {
			fmt.Println("[]")
		} else {
			fmt.Println("📭 No history found.")
		}
		return
	}

	switch historyFormat {
	case "json":
		printUnifiedJSON(records)
	case "table":
		printUnifiedTable(records)
	default:
		printUnifiedDefault(records)
	}
}

// collectUnifiedRecords gathers records from both tables based on --type filter.
func collectUnifiedRecords(database *db.DB) []unifiedRecord {
	var records []unifiedRecord

	// Include move_history records
	if historyType == "all" || historyType == "move" || historyType == "rename" {
		moves, err := database.ListMoveHistory(historyLimit)
		if err != nil {
			errlog.Warn("Error reading move history: %v", err)
		}
		for _, m := range moves {
			recType := "move"
			if m.FromPath != "" && m.ToPath != "" {
				// If from and to share the same directory, it's a rename
				if dirOf(m.FromPath) == dirOf(m.ToPath) {
					recType = "rename"
				}
			}
			// Apply type filter
			if historyType != "all" && historyType != recType {
				continue
			}
			detail := fmt.Sprintf("%s → %s", m.OriginalFileName, m.NewFileName)
			records = append(records, unifiedRecord{
				Source:    "move",
				ID:        m.ID,
				Type:      recType,
				Detail:    detail,
				FromPath:  m.FromPath,
				ToPath:    m.ToPath,
				Timestamp: m.MovedAt,
				Undone:    m.Undone,
			})
		}
	}

	// Include action_history records
	if shouldIncludeActions() {
		var actions []db.ActionRecord
		var err error

		switch historyType {
		case "scan":
			// Show both scan_add and scan_remove
			adds, _ := database.ListActionsByType(db.ActionScanAdd, historyLimit)
			removes, _ := database.ListActionsByType(db.ActionScanRemove, historyLimit)
			actions = append(adds, removes...)
		case "delete":
			actions, err = database.ListActionsByType(db.ActionDelete, historyLimit)
		case "popout":
			actions, err = database.ListActionsByType(db.ActionPopout, historyLimit)
		case "rescan":
			actions, err = database.ListActionsByType(db.ActionRescanUpdate, historyLimit)
		default: // "all"
			actions, err = database.ListActions(historyLimit)
		}
		if err != nil {
			errlog.Warn("Error reading action history: %v", err)
		}

		for _, a := range actions {
			detail := a.Detail
			if detail == "" {
				detail = string(a.ActionType)
			}
			records = append(records, unifiedRecord{
				Source:    "action",
				ID:        a.ID,
				Type:      string(a.ActionType),
				Detail:    detail,
				BatchID:   a.BatchID,
				Timestamp: a.CreatedAt,
				Undone:    a.Undone,
			})
		}
	}

	// Sort by timestamp descending (newest first)
	sortRecordsByTimestamp(records)

	// Apply --since filter
	if historySince != "" {
		var filtered []unifiedRecord
		for _, r := range records {
			if r.Timestamp >= historySince {
				filtered = append(filtered, r)
			}
		}
		records = filtered
	}

	// Apply limit
	if len(records) > historyLimit {
		records = records[:historyLimit]
	}

	return records
}

func shouldIncludeActions() bool {
	switch historyType {
	case "move", "rename":
		return false
	default:
		return true
	}
}

func showBatchHistory(database *db.DB) {
	actions, err := database.ListActionsByBatch(historyBatch)
	if err != nil {
		errlog.Error("Error reading batch %s: %v", historyBatch, err)
		return
	}
	if len(actions) == 0 {
		// Try partial match
		allActions, listErr := database.ListActions(200)
		if listErr != nil {
			errlog.Error("Error reading actions: %v", listErr)
			return
		}
		for _, a := range allActions {
			if len(a.BatchID) >= len(historyBatch) && a.BatchID[:len(historyBatch)] == historyBatch {
				actions = append(actions, a)
			}
		}
		if len(actions) == 0 {
			fmt.Printf("📭 No actions found for batch: %s\n", historyBatch)
			return
		}
	}

	fmt.Printf("📋 Batch: %s (%d actions)\n\n", historyBatch, len(actions))
	for _, a := range actions {
		status := "✅"
		if a.Undone {
			status = "↩️ "
		}
		detail := a.Detail
		if detail == "" {
			detail = string(a.ActionType)
		}
		fmt.Printf("  %s [%s] %s\n", status, a.ActionType, detail)
		fmt.Printf("     ID: %d  Created: %s\n\n", a.ID, a.CreatedAt)
	}
}

// ---------------------------------------------------------------------------
// Output formatters
// ---------------------------------------------------------------------------

func printUnifiedDefault(records []unifiedRecord) {
	fmt.Printf("📋 History (%d records)\n\n", len(records))

	for _, r := range records {
		status := "✅"
		if r.Undone {
			status = "↩️ "
		}

		icon := typeIcon(r.Type)
		fmt.Printf("  %s %s %-14s  %s\n", status, icon, r.Type, r.Timestamp)
		fmt.Printf("     %s\n", r.Detail)

		if r.FromPath != "" {
			fmt.Printf("     Path: %s → %s\n", r.FromPath, r.ToPath)
		}
		if r.BatchID != "" {
			fmt.Printf("     Batch: %s\n", r.BatchID[:minInt(8, len(r.BatchID))])
		}
		fmt.Println()
	}
}

func printUnifiedJSON(records []unifiedRecord) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(records); err != nil {
		errlog.Error("JSON encode error: %v", err)
	}
}

func printUnifiedTable(records []unifiedRecord) {
	printHistoryTableUnified(records)
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func typeIcon(t string) string {
	switch t {
	case "move":
		return "📁"
	case "rename":
		return "✏️ "
	case "scan_add":
		return "➕"
	case "scan_remove":
		return "➖"
	case "delete":
		return "🗑 "
	case "popout":
		return "📤"
	case "restore":
		return "♻️ "
	case "rescan_update":
		return "🔄"
	default:
		return "📋"
	}
}

func dirOf(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' || path[i] == '\\' {
			return path[:i]
		}
	}
	return ""
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// sortRecordsByTimestamp sorts unified records by timestamp descending.
func sortRecordsByTimestamp(records []unifiedRecord) {
	for i := 1; i < len(records); i++ {
		for j := i; j > 0 && records[j].Timestamp > records[j-1].Timestamp; j-- {
			records[j], records[j-1] = records[j-1], records[j]
		}
	}
}
