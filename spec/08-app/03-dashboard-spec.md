# Dashboard Specification — Media Library Dashboard

> **Version**: 1.0  
> **Created**: 09-Apr-2026  
> **Status**: Approved for phased implementation  
> **Target Stack**: React 18 + Tailwind CSS (existing Lovable project)

---

## 1. Purpose

Build a web-based dashboard that reads media records from the SQLite database populated by `movie scan` (and other ingestion commands), and renders them as a browsable, filterable, card-based interface.

---

## 2. Data Source & Feasibility

### 2.1 Available Fields (from `media` table)

| Field | DB Column | Available | Source |
|-------|-----------|-----------|--------|
| Title | `title` | ✅ Always | Filename cleaner or TMDb |
| Description | `description` | ✅ When TMDb key set | TMDb API `overview` |
| Poster/Image | `thumbnail_path` | ✅ When TMDb key set | Downloaded JPG from TMDb |
| TMDb Rating | `tmdb_rating` | ✅ When TMDb key set | TMDb `vote_average` |
| IMDb Rating | `imdb_rating` | ✅ When TMDb key set | TMDb movie details |
| Genres | `genre` | ✅ When TMDb key set | Comma-separated genre names |
| Type (Category) | `type` | ✅ Always | `"movie"` or `"tv"` — auto-detected |
| Director | `director` | ✅ When TMDb key set | TMDb credits (Job = "Director") |
| Cast/Actors | `cast_list` | ✅ When TMDb key set | Top 10 cast from TMDb credits |
| Year | `year` | ✅ Usually | Extracted from filename or TMDb |
| Popularity | `popularity` | ✅ When TMDb key set | TMDb popularity score |
| File Size | `file_size` | ✅ Always | `os.Stat()` during scan |
| File Path | `current_file_path` | ✅ Always | Tracked across moves |
| File Extension | `file_extension` | ✅ Always | Parsed from filename |
| Tags | `tags` table | ✅ Optional | User-assigned via `movie tag` |

### 2.2 Feasibility Verdict

**All requested dashboard fields are available.** The `media` table already stores title, description, poster path, ratings, genres, type (category), director, and cast. The only requirement is that the user has a TMDb API key configured — without it, only title/year/type/file info are available.

### 2.3 Limitations

- **Poster images** are stored as local file paths (`thumbnail_path`). The dashboard must either serve them from the filesystem or read them via a data URL / static serve.
- **Actor details** are limited to names (no photos, no character names in DB). Character info exists in TMDb credits but is not currently persisted.
- **Genres** are stored as a comma-separated string, not normalized. Filtering requires string splitting.

---

## 3. Data Model for Dashboard

### 3.1 MediaItem (Frontend Interface)

```typescript
interface MediaItem {
  id: number;
  title: string;
  cleanTitle: string;
  year: number;
  type: "movie" | "tv";
  tmdbId: number;
  imdbId: string;
  description: string;
  imdbRating: number;
  tmdbRating: number;
  popularity: number;
  genres: string[];        // split from comma-separated genre string
  director: string;
  cast: string[];          // split from comma-separated cast_list
  thumbnailPath: string;   // local filesystem path to poster JPG
  thumbnailUrl: string;    // resolved URL for <img> display
  currentFilePath: string;
  fileSize: number;        // bytes
  fileExtension: string;
  tags: string[];          // from tags table
  metadataStatus: "full" | "partial" | "filename-only";
}
```

### 3.2 Metadata Status Logic

| Status | Condition |
|--------|-----------|
| `full` | `tmdb_id > 0` AND `description != ""` AND `thumbnail_path != ""` |
| `partial` | `tmdb_id > 0` but missing description or poster |
| `filename-only` | `tmdb_id == 0` — scanned without TMDb enrichment |

---

## 4. Command Behavior Reference

### 4.1 `movie scan <folder>`

| Aspect | Detail |
|--------|--------|
| **Input** | Folder path (arg) or `scan_dir` config fallback |
| **Validation** | Checks folder exists, is a directory |
| **Discovery** | Reads top-level entries; for directories, looks for first video file inside |
| **File Filter** | `cleaner.IsVideoFile()` — matches `.mkv`, `.mp4`, `.avi`, `.mov`, `.wmv`, `.flv`, `.webm` |
| **Cleaning** | `cleaner.Clean(filename)` → removes junk (resolution, codec, release group), extracts year, detects TV patterns |
| **Dedup** | Checks if `original_file_path` already exists in DB → skips |
| **TMDb Fetch** | If API key set: `SearchMulti(title)` → picks best match → `GetMovieDetails` or `GetTVDetails` → `GetCredits` → `DownloadPoster` |
| **Storage** | `InsertMedia()` or `UpdateMediaByTmdbID()` on duplicate |
| **History** | `InsertScanHistory(folder, totalFiles, movies, tvShows)` |
| **Output** | Per-file progress + summary counts |
| **Error: Empty folder** | Prints "No video files found" if no video files detected |
| **Error: Invalid folder** | Prints "Folder not found" and exits |
| **Error: No API key** | Warns but proceeds — stores filename-only records |
| **Error: TMDb no match** | Warns per file, stores record without metadata |

### 4.2 `movie search <query>`

| Aspect | Detail |
|--------|--------|
| **Input** | One or more search terms |
| **Requires** | TMDb API key (exits if missing) |
| **Flow** | `SearchMulti(query)` → show up to 15 results → user picks number → fetch full details → save to DB |
| **Does NOT require** | File to exist in library — pure metadata lookup |
| **Downstream** | Record available in `ls`, `info`, `stats` |

### 4.3 `movie info <id-or-title>`

| Aspect | Detail |
|--------|--------|
| **Input** | Numeric DB ID or title string |
| **Resolution** | `resolveMediaByQuery()`: ID lookup → exact title → prefix match → first fuzzy result |
| **Fallback** | If not in DB → TMDb search → auto-save first result → display |
| **Output** | Full detail view (title, year, type, ratings, genre, director, cast, description, paths) |

### 4.4 `movie ls`

| Aspect | Detail |
|--------|--------|
| **Input** | None |
| **Flow** | Paginated list from DB (configurable page size) |
| **Navigation** | `N` next, `P` previous, `Q` quit, number for detail view |
| **Shows** | Clean title, year, rating, type icon |

### 4.5 `movie move [directory]`

| Aspect | Detail |
|--------|--------|
| **Input** | Optional source directory |
| **Flow** | Browse video files → select → choose destination → confirm → move |
| **Cross-drive** | `MoveFile()` with `os.Rename` → EXDEV fallback → `io.Copy` + `os.Remove` |
| **Tracking** | Updates `current_file_path` + inserts `move_history` for undo |

### 4.6 `movie rename`

| Aspect | Detail |
|--------|--------|
| **Input** | None (operates on all DB records) |
| **Flow** | Find files where current name ≠ clean name → show list → confirm → batch rename |
| **Undo** | Each rename tracked in `move_history` |

### 4.7 `movie undo`

| Aspect | Detail |
|--------|--------|
| **Input** | None |
| **Flow** | Get last undone move → show from/to → confirm `[y/N]` → revert |
| **Safety** | Checks file exists at destination before reverting |

### 4.8 `movie stats`

| Aspect | Detail |
|--------|--------|
| **Output** | Total movies/TV/all, file size stats (total/largest/smallest/avg), top 10 genres with bar chart, average ratings |

### 4.9 `movie suggest [N]`

| Aspect | Detail |
|--------|--------|
| **Input** | Optional count (default 10) |
| **Flow** | Choose Movie/TV/Random → analyze top genres → fetch TMDb recommendations from random library items → fill gaps with trending |

### 4.10 `movie tag add|remove|list`

| Aspect | Detail |
|--------|--------|
| **Subcommands** | `add <id> <tag>`, `remove <id> <tag>`, `list [id]` |
| **Storage** | `tags` table with unique `(media_id, tag)` constraint |

### 4.11 `movie config [get|set]`

| Aspect | Detail |
|--------|--------|
| **Keys** | `movies_dir`, `tv_dir`, `archive_dir`, `scan_dir`, `tmdb_api_key`, `page_size` |
| **Display** | API key is masked (first 4 + ... + last 4) |

### 4.12 `movie play <id>`

| Aspect | Detail |
|--------|--------|
| **Input** | Numeric media ID |
| **Flow** | Look up file path → verify exists → `open`/`xdg-open`/`cmd /c start` |

---

## 5. Dashboard UI Specification

### 5.1 Layout Structure

```
+------------------------------------------------------------------+
| HEADER                                                            |
|  Dashboard Title | Source Path | Summary Counts                  |
+------------------------------------------------------------------+
| FILTERS                                                           |
|  [Search by title...] [Genre ▾] [Type: All|Movie|TV] [Sort ▾]   |
+------------------------------------------------------------------+
| CONTENT GRID                                                      |
|  +--------+  +--------+  +--------+  +--------+                  |
|  | Poster |  | Poster |  | Poster |  | Poster |                  |
|  | Title  |  | Title  |  | Title  |  | Title  |                  |
|  | Rating |  | Rating |  | Rating |  | Rating |                  |
|  | Genres |  | Genres |  | Genres |  | Genres |                  |
|  | Cast   |  | Cast   |  | Cast   |  | Cast   |                  |
|  +--------+  +--------+  +--------+  +--------+                  |
+------------------------------------------------------------------+
| DETAIL MODAL (on card click)                                      |
|  Large Poster | Full Description | Director | Cast | File Info   |
+------------------------------------------------------------------+
```

### 5.2 Component Breakdown

#### Phase 1: Data Layer + Core Layout
- `MediaDataProvider` — context provider with mock data (JSON)
- `DashboardLayout` — header + content area shell
- `DashboardHeader` — title, library stats (total movies, total TV, total size)

#### Phase 2: Media Card Grid
- `MediaCard` — poster, title, year, rating stars, genre chips, cast preview (first 3)
- `MediaGrid` — responsive grid of `MediaCard` components
- Fallback UI for missing poster (placeholder with title initial)
- Fallback text for missing description/rating/cast

#### Phase 3: Filters + Search
- `SearchBar` — title search with debounce
- `GenreFilter` — multi-select dropdown extracted from all unique genres
- `TypeFilter` — toggle: All / Movie / TV
- `SortSelect` — sort by: Title A-Z, Rating High-Low, Year New-Old, Popularity

#### Phase 4: Detail Modal
- `MediaDetailModal` — full-screen or dialog overlay
- Large poster image
- Full description (no truncation)
- Director + full cast list
- Ratings (TMDb + IMDb if available)
- All genres as tags
- File info: path, size, extension
- Tags display
- Metadata status indicator

#### Phase 5: Statistics Panel
- `StatsPanel` — summary cards + genre distribution chart
- Total counts (movies, TV, all)
- File size stats (total, average, largest)
- Top genres bar chart (using recharts — already installed)
- Average ratings

### 5.3 UI Behavior Rules

| Rule | Behavior |
|------|----------|
| Missing poster | Show gradient placeholder with first letter of title |
| Long description | Truncate to 120 chars on card, full in modal |
| Missing rating | Show "N/A" in muted text |
| Missing cast | Hide cast section on card, show "No cast info" in modal |
| Missing genre | Hide genre chips |
| Genres display | Render as colored tag chips (max 3 on card, all in modal) |
| Rating display | Star icon + numeric value, colored by quality (green ≥7, yellow ≥5, red <5) |
| Type badge | Movie = blue badge, TV = purple badge |
| Responsive | 1 col mobile, 2 col tablet, 3-4 col desktop |
| Empty state | "No media found. Run `movie scan` to get started." |

### 5.4 Theme

Use existing design system tokens from `index.css` and `tailwind.config.ts`. Card-based layout with:
- Rounded corners (`rounded-lg`)
- Subtle shadow (`shadow-md`)
- Background cards on `card` / `card-foreground` tokens
- Consistent spacing (`gap-4`, `p-4`)
- Genre chips: muted background with border

---

## 6. Data Loading Strategy

Since this is a Lovable web project and the Go CLI writes to a local SQLite file, the dashboard needs a data bridge. Options:

### 6.1 Option A: Static JSON Export (Recommended for Phase 1)
- Add a `movie export` command that dumps the `media` table as JSON to a known location
- Dashboard loads this JSON file
- Simple, no backend needed, works with existing infrastructure

### 6.2 Option B: REST API (Future)
- Add a `movie serve` command that starts a local HTTP server
- Exposes `/api/media`, `/api/stats`, `/api/genres` endpoints
- Dashboard fetches from localhost

### 6.3 Option C: Mock Data (Development)
- Start with hardcoded mock data matching the `MediaItem` interface
- Build all UI components against mock data
- Replace with real data source later

**Decision: Start with Option C (mock data), then implement Option A (JSON export) when UI is complete.**

---

## 7. Implementation Phases

| Phase | Scope | Deliverables |
|-------|-------|-------------|
| **Phase 1** | Data + Layout | `MediaItem` type, mock data (10 items), `DashboardLayout`, `DashboardHeader` with stats |
| **Phase 2** | Card Grid | `MediaCard`, `MediaGrid`, responsive layout, poster fallback, rating colors |
| **Phase 3** | Filters + Search | `SearchBar`, `GenreFilter`, `TypeFilter`, `SortSelect`, filtered state |
| **Phase 4** | Detail Modal | `MediaDetailModal` with full info, click-to-open from card |
| **Phase 5** | Stats Panel | `StatsPanel` with recharts genre chart, summary cards, rating averages |

Each phase is self-contained and can be implemented independently when the user says "next."

---

## 8. Acceptance Criteria

### Phase 1
- GIVEN the dashboard loads WHEN there are 10 mock media items THEN the header shows "10 titles, X movies, Y TV shows"
- GIVEN the layout renders THEN it uses semantic design tokens, not hardcoded colors

### Phase 2
- GIVEN a media item with a poster WHEN the card renders THEN the poster image is displayed
- GIVEN a media item without a poster WHEN the card renders THEN a gradient placeholder with the title initial is shown
- GIVEN a rating ≥ 7.0 WHEN displayed THEN it uses green color
- GIVEN genres exist WHEN the card renders THEN max 3 genre chips are shown

### Phase 3
- GIVEN the user types "dark" in search WHEN filtering THEN only titles containing "dark" are shown
- GIVEN "Action" is selected in genre filter WHEN filtering THEN only items with Action genre appear
- GIVEN "TV" is selected in type filter THEN only TV shows are shown

### Phase 4
- GIVEN the user clicks a card THEN a detail modal opens with full metadata
- GIVEN the modal is open WHEN Escape is pressed THEN the modal closes

### Phase 5
- GIVEN media records exist WHEN stats panel renders THEN genre bar chart shows top genres
- GIVEN file sizes exist THEN total, average, and largest are displayed in human-readable format
