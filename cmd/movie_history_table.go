// movie_history_table.go — table-formatted output for unified history
package cmd

import (
	"fmt"
	"strings"
)

func printHistoryTableUnified(records []unifiedRecord) {
	idW := 6
	typeW := 14
	statusW := 6
	dateW := 19
	detailW := 40

	fmt.Println()
	fmt.Printf("  %-*s │ %-*s │ %-*s │ %-*s │ %-*s\n",
		idW, "ID",
		typeW, "Type",
		statusW, "Status",
		dateW, "Date",
		detailW, "Detail")

	fmt.Printf("  %s─┼─%s─┼─%s─┼─%s─┼─%s\n",
		strings.Repeat("─", idW),
		strings.Repeat("─", typeW),
		strings.Repeat("─", statusW),
		strings.Repeat("─", dateW),
		strings.Repeat("─", detailW))

	for _, r := range records {
		status := "OK"
		if r.Undone {
			status = "Undone"
		}

		prefix := r.Source[0:1] // "m" or "a"
		idStr := fmt.Sprintf("%s-%d", prefix, r.ID)

		fmt.Printf("  %-*s │ %-*s │ %-*s │ %-*s │ %-*s\n",
			idW, idStr,
			typeW, truncate(r.Type, typeW),
			statusW, status,
			dateW, truncate(r.Timestamp, dateW),
			detailW, truncate(r.Detail, detailW))
	}

	fmt.Printf("  %s─┴─%s─┴─%s─┴─%s─┴─%s\n",
		strings.Repeat("─", idW),
		strings.Repeat("─", typeW),
		strings.Repeat("─", statusW),
		strings.Repeat("─", dateW),
		strings.Repeat("─", detailW))

	fmt.Printf("\n  Total: %d records\n\n", len(records))
}

// printHistoryTable is kept for backward compatibility (move-only table).
func printHistoryTable(records []db.MoveRecord) {
	var unified []unifiedRecord
	for _, m := range records {
		recType := "move"
		if dirOf(m.FromPath) == dirOf(m.ToPath) {
			recType = "rename"
		}
		unified = append(unified, unifiedRecord{
			Source:    "move",
			ID:        m.ID,
			Type:      recType,
			Detail:    fmt.Sprintf("%s → %s", m.OriginalFileName, m.NewFileName),
			FromPath:  m.FromPath,
			ToPath:    m.ToPath,
			Timestamp: m.MovedAt,
			Undone:    m.Undone,
		})
	}
	printHistoryTableUnified(unified)
}
