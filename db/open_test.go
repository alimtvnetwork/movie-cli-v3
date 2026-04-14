package db

import (
	"testing"
)

// ─── Schema / Migration ─────────────────────────────────────

func TestMigrate(t *testing.T) {
	d := openTestDB(t)
	// Verify tables exist by querying them
	for _, table := range []string{"media", "move_history", "config", "scan_history", "tags", "watchlist"} {
		var name string
		err := d.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name=?", table).Scan(&name)
		if err != nil {
			t.Errorf("table %q not found: %v", table, err)
		}
	}
}

func TestDefaultConfig(t *testing.T) {
	d := openTestDB(t)
	val, err := d.GetConfig("movies_dir")
	if err != nil {
		t.Fatalf("get config: %v", err)
	}
	if val != "~/Movies" {
		t.Errorf("movies_dir = %q, want ~/Movies", val)
	}
}

// ─── Media CRUD ─────────────────────────────────────────────

func TestInsertAndGetMedia(t *testing.T) {
	d := openTestDB(t)
	id := seedMedia(t, d, "Inception", 27205)

	m, err := d.GetMediaByID(id)
	if err != nil {
		t.Fatalf("get by id: %v", err)
	}
	if m.Title != "Inception" || m.TmdbID != 27205 {
		t.Errorf("got title=%q tmdb=%d", m.Title, m.TmdbID)
	}

	m2, err := d.GetMediaByTmdbID(27205)
	if err != nil {
		t.Fatalf("get by tmdb: %v", err)
	}
	if m2.ID != id {
		t.Errorf("tmdb lookup returned id %d, want %d", m2.ID, id)
	}
}

func TestInsertMediaAllowsMissingTMDbID(t *testing.T) {
	d := openTestDB(t)

	if _, err := d.InsertMedia(&Media{Title: "Local One", CleanTitle: "local one", Type: "movie"}); err != nil {
		t.Fatalf("first insert: %v", err)
	}
	if _, err := d.InsertMedia(&Media{Title: "Local Two", CleanTitle: "local two", Type: "movie"}); err != nil {
		t.Fatalf("second insert with empty tmdb_id should store NULL: %v", err)
	}
}

func TestUpdateMediaByTmdbID(t *testing.T) {
	d := openTestDB(t)
	seedMedia(t, d, "Inception", 27205)

	err := d.UpdateMediaByTmdbID(&Media{
		Title:      "Inception (Updated)",
		CleanTitle: "inception updated",
		Year:       2010,
		Type:       "movie",
		TmdbID:     27205,
	})
	if err != nil {
		t.Fatalf("update: %v", err)
	}

	m, _ := d.GetMediaByTmdbID(27205)
	if m.Title != "Inception (Updated)" {
		t.Errorf("title = %q after update", m.Title)
	}
}

func TestUpdateMediaPath(t *testing.T) {
	d := openTestDB(t)
	id := seedMedia(t, d, "Test", 1)

	if err := d.UpdateMediaPath(id, "/new/path.mkv"); err != nil {
		t.Fatalf("update path: %v", err)
	}
	m, _ := d.GetMediaByID(id)
	if m.CurrentFilePath != "/new/path.mkv" {
		t.Errorf("path = %q", m.CurrentFilePath)
	}
}

func TestListMedia(t *testing.T) {
	d := openTestDB(t)
	seedMedia(t, d, "Alpha", 1)
	seedMedia(t, d, "Beta", 2)
	seedMedia(t, d, "Gamma", 3)

	list, err := d.ListMedia(0, 2)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(list) != 2 {
		t.Errorf("got %d items, want 2", len(list))
	}

	// Second page
	list2, _ := d.ListMedia(2, 10)
	if len(list2) != 1 {
		t.Errorf("page 2: got %d items, want 1", len(list2))
	}
}

func TestSearchMedia(t *testing.T) {
	d := openTestDB(t)
	seedMedia(t, d, "The Matrix", 603)
	seedMedia(t, d, "Inception", 27205)

	results, err := d.SearchMedia("matrix")
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	if len(results) != 1 || results[0].Title != "The Matrix" {
		t.Errorf("search returned %d results", len(results))
	}
}

func TestCountMedia(t *testing.T) {
	d := openTestDB(t)
	seedMedia(t, d, "Movie1", 1)
	seedMedia(t, d, "Movie2", 2)

	count, err := d.CountMedia("")
	if err != nil {
		t.Fatalf("count: %v", err)
	}
	if count != 2 {
		t.Errorf("count = %d, want 2", count)
	}

	countMovies, _ := d.CountMedia("movie")
	if countMovies != 2 {
		t.Errorf("movie count = %d", countMovies)
	}

	countTV, _ := d.CountMedia("tv")
	if countTV != 0 {
		t.Errorf("tv count = %d", countTV)
	}
}

func TestDeleteMedia(t *testing.T) {
	d := openTestDB(t)
	id := seedMedia(t, d, "ToDelete", 999)

	if err := d.DeleteMedia(id); err != nil {
		t.Fatalf("delete: %v", err)
	}
	_, err := d.GetMediaByID(id)
	if err == nil {
		t.Error("expected error after delete")
	}
}

// ─── Config ─────────────────────────────────────────────────

func TestConfigSetGet(t *testing.T) {
	d := openTestDB(t)

	if err := d.SetConfig("test_key", "test_value"); err != nil {
		t.Fatalf("set: %v", err)
	}
	val, err := d.GetConfig("test_key")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if val != "test_value" {
		t.Errorf("config = %q", val)
	}

	// Overwrite
	d.SetConfig("test_key", "new_value")
	val2, _ := d.GetConfig("test_key")
	if val2 != "new_value" {
		t.Errorf("overwrite = %q", val2)
	}
}

// ─── Tags ───────────────────────────────────────────────────

func TestTags(t *testing.T) {
	d := openTestDB(t)
	id := seedMedia(t, d, "Tagged Movie", 100)

	if err := d.AddTag(int(id), "favorite"); err != nil {
		t.Fatalf("add tag: %v", err)
	}
	if err := d.AddTag(int(id), "action"); err != nil {
		t.Fatalf("add tag: %v", err)
	}

	tags, err := d.GetTagsByMediaID(int(id))
	if err != nil {
		t.Fatalf("get tags: %v", err)
	}
	if len(tags) != 2 {
		t.Errorf("got %d tags, want 2", len(tags))
	}

	// Duplicate tag should fail
	if err := d.AddTag(int(id), "favorite"); err == nil {
		t.Error("expected error on duplicate tag")
	}

	// Remove tag
	ok, err := d.RemoveTag(int(id), "action")
	if err != nil || !ok {
		t.Errorf("remove: ok=%v err=%v", ok, err)
	}

	// Remove non-existent
	ok2, _ := d.RemoveTag(int(id), "nonexistent")
	if ok2 {
		t.Error("expected false for non-existent tag")
	}

	// Tag counts
	counts, err := d.GetAllTagCounts()
	if err != nil {
		t.Fatalf("tag counts: %v", err)
	}
	if len(counts) != 1 || counts[0].Tag != "favorite" {
		t.Errorf("counts = %+v", counts)
	}
}

// ─── Move History ───────────────────────────────────────────

func TestMoveHistory(t *testing.T) {
	d := openTestDB(t)
	id := seedMedia(t, d, "Moved Movie", 200)

	err := d.InsertMoveHistory(id, "/old/path.mkv", "/new/path.mkv", "old.mkv", "new.mkv")
	if err != nil {
		t.Fatalf("insert history: %v", err)
	}

	rec, err := d.GetLastMove()
	if err != nil {
		t.Fatalf("get last move: %v", err)
	}
	if rec.FromPath != "/old/path.mkv" || rec.ToPath != "/new/path.mkv" {
		t.Errorf("move record: from=%q to=%q", rec.FromPath, rec.ToPath)
	}

	if err := d.MarkMoveUndone(rec.ID); err != nil {
		t.Fatalf("undo: %v", err)
	}

	// No un-undone moves left
	_, err = d.GetLastMove()
	if err == nil {
		t.Error("expected no more un-undone moves")
	}
}

// ─── Scan History ───────────────────────────────────────────

func TestScanHistory(t *testing.T) {
	d := openTestDB(t)
	err := d.InsertScanHistory("/downloads", 50, 30, 20)
	if err != nil {
		t.Fatalf("insert scan history: %v", err)
	}
}

// ─── Watchlist ──────────────────────────────────────────────

func TestWatchlist(t *testing.T) {
	d := openTestDB(t)

	// Add to watchlist
	if err := d.AddToWatchlist(550, "Fight Club", 1999, "movie", 0); err != nil {
		t.Fatalf("add: %v", err)
	}

	// List to-watch
	list, err := d.ListWatchlist("to-watch")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(list) != 1 || list[0].Title != "Fight Club" {
		t.Errorf("watchlist = %+v", list)
	}

	// Mark watched
	if err := d.MarkWatched(550); err != nil {
		t.Fatalf("mark watched: %v", err)
	}
	entry, _ := d.GetWatchlistByTmdbID(550)
	if entry.Status != "watched" {
		t.Errorf("status = %q", entry.Status)
	}

	// Undo (back to to-watch)
	d.MarkToWatch(550)
	entry2, _ := d.GetWatchlistByTmdbID(550)
	if entry2.Status != "to-watch" {
		t.Errorf("undo status = %q", entry2.Status)
	}

	// List all
	all, _ := d.ListWatchlist("")
	if len(all) != 1 {
		t.Errorf("all = %d", len(all))
	}

	// Remove
	d.RemoveFromWatchlist(550)
	after, _ := d.ListWatchlist("")
	if len(after) != 0 {
		t.Errorf("after remove = %d", len(after))
	}
}

// ─── Watchlist upsert ───────────────────────────────────────

func TestWatchlistUpsert(t *testing.T) {
	d := openTestDB(t)
	d.AddToWatchlist(550, "Fight Club", 1999, "movie", 0)
	// Re-add should upsert, not error
	if err := d.AddToWatchlist(550, "Fight Club (Updated)", 1999, "movie", 0); err != nil {
		t.Fatalf("upsert: %v", err)
	}
	entry, _ := d.GetWatchlistByTmdbID(550)
	if entry.Title != "Fight Club (Updated)" {
		t.Errorf("upsert title = %q", entry.Title)
	}
}

// ─── File Size Stats ────────────────────────────────────────

func TestFileSizeStats(t *testing.T) {
	d := openTestDB(t)
	seedMedia(t, d, "Small", 1) // 700MB from seedMedia
	seedMedia(t, d, "Big", 2)

	total, largest, smallest, err := d.FileSizeStats()
	if err != nil {
		t.Fatalf("stats: %v", err)
	}
	if total != 700*1024*1024*2 {
		t.Errorf("total = %d", total)
	}
	if largest != smallest {
		t.Errorf("largest=%d smallest=%d (should be equal)", largest, smallest)
	}
}

// ─── Top Genres ─────────────────────────────────────────────

func TestTopGenres(t *testing.T) {
	d := openTestDB(t)
	id1, _ := d.InsertMedia(&Media{
		Title: "A", CleanTitle: "a", Year: 2024, Type: "movie", TmdbID: 1,
		Genre: "Action, Drama", OriginalFilePath: "/a.mkv",
	})
	_ = id1
	d.InsertMedia(&Media{
		Title: "B", CleanTitle: "b", Year: 2024, Type: "movie", TmdbID: 2,
		Genre: "Action, Comedy", OriginalFilePath: "/b.mkv",
	})

	genres, err := d.TopGenres(10)
	if err != nil {
		t.Fatalf("genres: %v", err)
	}
	if genres["Action"] != 2 {
		t.Errorf("Action = %d, want 2", genres["Action"])
	}
	if genres["Drama"] != 1 {
		t.Errorf("Drama = %d", genres["Drama"])
	}
}
