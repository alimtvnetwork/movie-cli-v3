# 10 — CLI Output Format Specification

> Version: 1.1 | Updated: 2026-04-14

## Overview

This spec defines the visual output format for the `movie scan` command.
The format is inspired by [gitmap](https://github.com/alimtvnetwork/gitmap-v2)
and provides a structured, human-readable terminal experience with numbered
items, categorized icons, ratings, and a tree-style output file listing.

This document is self-contained for AI handoff — no code reading required.

---

## 1. Header Block

A double-line box banner showing the version is printed at the top:

```
  ╔══════════════════════════════════════╗
  ║       🎬  Movie CLI v1.31.0         ║
  ╚══════════════════════════════════════╝

  📂 Scanning: /path/to/folder
  🔄 Mode: recursive (all subdirectories)
  📁 Output: /path/to/folder/.movie-output
```

### Rules
- Box width: fixed at 38 chars inner
- Version label is centered: `🎬  Movie CLI {version}`
- Mode line only appears if `--recursive` is used
- Output line is hidden in `--dry-run` mode
- All lines indented with 2 spaces

---

## 2. Flag Reference

| Flag            | Short | Type   | Default   | Description                                       |
|-----------------|-------|--------|-----------|---------------------------------------------------|
| `--recursive`   | `-r`  | bool   | `false`   | Scan all subdirectories recursively                |
| `--depth`       | `-d`  | int    | `0`       | Max subdirectory depth (0 = unlimited)             |
| `--dry-run`     |       | bool   | `false`   | Preview without writing to DB or output            |
| `--format`      |       | string | `default` | Output format: `default`, `table`, or `json`       |
| `--rest`        |       | bool   | `false`   | Start REST server + open browser after scan        |
| `--port`        |       | int    | `8086`    | Port for REST server when using `--rest`           |

### Flag Interactions & Edge Cases

| Combination               | Behavior                                                    |
|---------------------------|-------------------------------------------------------------|
| `--dry-run --rest`        | Dry run executes normally. `--rest` is ignored (no output). |
| `--format json --rest`    | JSON printed to stdout first, then REST server starts.      |
| `--format table --rest`   | Table printed first, then REST server starts.               |
| `--format json --dry-run` | JSON printed to stdout. No files written. No REST.          |
| `--depth` without `-r`    | `--depth` is silently ignored.                              |

---

## 3. Scanned Items Section

```
  ■ Scanned Items
  ──────────────────────────────────────────
```

Each item is printed as a numbered entry with an icon indicating type:

```
  1. 🎬 Inception (2010) [movie]
     └─ Inception.2010.1080p.BluRay.x264.mkv
     ⭐ 8.4  Inception

  2. 📺 Breaking Bad (2008) [tv]
     └─ Breaking.Bad.S01E01.720p.mkv
     ⭐ 9.5  Breaking Bad

  3. 🎬 The Dark Knight (2008) [movie]
     └─ The.Dark.Knight.2008.REMUX.mkv
     ⏩ Already in database, skipping
```

### Item Format

```
  {index}. {icon} {clean_title} ({year}) [{type}]
     └─ {original_filename}
     ⭐ {rating}  {tmdb_title}
```

### Icons
| Type    | Icon |
|---------|------|
| Movie   | 🎬   |
| TV Show | 📺   |

### Metadata Lines (after `└─`)

These lines appear conditionally after the filename, in this order:

| Priority | Condition         | Line                                 |
|----------|-------------------|--------------------------------------|
| 1        | Already in DB     | `⏩ Already in database, skipping`    |
| 2        | TMDb matched      | `⭐ {rating}  {title}`               |
| 3        | Thumbnail saved   | `🖼️  Thumbnail saved`               |
| 4        | TMDb warning      | `⚠️  no TMDb match for '{query}'`    |

If an item is skipped (already in DB), lines 2-4 are not printed.

---

## 4. Summary Section

```
  ■ Summary
  ──────────────────────────────────────────
  📊 Scan Complete!
     Total files: 25
     Movies:      18
     TV Shows:    7
     Skipped:     3 (already in DB)
```

### Rules
- "Skipped" line only appears if `skipped > 0`
- In dry-run mode, title changes to "Dry Run Complete!"
- In dry-run mode, a tip line is appended:
  `💡 Run without --dry-run to actually scan and save.`

---

## 5. Output Files Section

Only shown when NOT in dry-run mode:

```
  ■ Output Files
  ──────────────────────────────────────────
  📁 /path/to/.movie-output/
  ├── 📄 summary.json      Scan report with metadata
  ├── 🌐 report.html       Interactive HTML report
  ├── 📁 json/movie/       Per-movie JSON metadata
  ├── 📁 json/tv/          Per-show JSON metadata
  └── 📁 thumbnails/       Movie poster thumbnails
```

### File Descriptions
| File/Dir         | Description                      |
|------------------|----------------------------------|
| `summary.json`   | Full scan report with counts     |
| `report.html`    | Interactive filterable HTML UI   |
| `json/movie/`    | One JSON file per movie          |
| `json/tv/`       | One JSON file per TV show        |
| `thumbnails/`    | Poster images: `{slug}-{id}.jpg` |

---

## 6. REST Server Section (--rest flag)

When `--rest` is specified, after the scan completes:

```
  🚀 Starting REST server on http://localhost:8086 ...
```

The server starts, the default browser opens the HTML report, and the CLI
blocks until Ctrl+C is pressed. All REST requests are logged to the error
log system (see Section 11).

---

## 7. Thumbnail Naming Convention

Thumbnails are saved to `.movie-output/thumbnails/` with the format:

```
{slug}-{tmdb_id}.jpg
```

Where `slug` = `ToSlug(clean_title)` + optional `-{year}`.

Examples:
- `inception-2010-27205.jpg`
- `breaking-bad-2008-1396.jpg`
- `the-dark-knight-2008-155.jpg`

Thumbnails are also copied to the database data directory (`data/thumbnails/`)
so the REST server can serve them at `/thumbnails/{filename}`.

---

## 8. JSON Output Mode (`--format json`)

When `--format json` is used, NO visual output is printed. Instead, a
single JSON object is written to stdout:

### Top-Level Schema

```json
{
  "scanned_folder": "/path/to/folder",
  "scanned_at": "2026-04-14T12:00:00Z",
  "dry_run": false,
  "total_files": 25,
  "movies": 18,
  "tv_shows": 7,
  "skipped": 3,
  "items": [...]
}
```

| Field            | Type     | Description                           |
|------------------|----------|---------------------------------------|
| `scanned_folder` | string   | Absolute path of scanned directory    |
| `scanned_at`     | string   | RFC 3339 timestamp (UTC)              |
| `dry_run`        | bool     | Whether this was a dry run            |
| `total_files`    | int      | Total video files found               |
| `movies`         | int      | Count of items classified as movie    |
| `tv_shows`       | int      | Count of items classified as TV       |
| `skipped`        | int      | Count of items already in database    |
| `items`          | array    | Array of `scanJSONItem` (see below)   |

### Item Schema (`scanJSONItem`)

```json
{
  "file_name": "Inception.2010.1080p.BluRay.x264.mkv",
  "file_path": "/movies/Inception.2010.1080p.BluRay.x264.mkv",
  "clean_title": "Inception",
  "year": 2010,
  "type": "movie",
  "tmdb_id": 27205,
  "tmdb_rating": 8.4,
  "genre": "Action, Science Fiction, Adventure",
  "status": "new"
}
```

| Field         | Type   | Omit If  | Description                              |
|---------------|--------|----------|------------------------------------------|
| `file_name`   | string | never    | Original filename on disk                |
| `file_path`   | string | never    | Full absolute path                       |
| `clean_title` | string | never    | Cleaned title after filename parsing     |
| `year`        | int    | `0`      | Extracted year from filename             |
| `type`        | string | never    | `"movie"` or `"tv"`                      |
| `tmdb_id`     | int    | `0`      | TMDb ID (0 if no match)                  |
| `tmdb_rating` | float  | `0`      | TMDb vote average (0 if no match)        |
| `genre`       | string | empty    | Comma-separated genre names              |
| `status`      | string | never    | `"new"`, `"skipped"`, or `"updated"`     |

---

## 9. Table Output Mode (`--format table`)

A fixed-width columnar output using Unicode box-drawing characters:

```
  #    │ File Name                      │ Clean Title                    │ Year  │ Type   │ Rating │ Status
  ─────┼────────────────────────────────┼────────────────────────────────┼───────┼────────┼────────┼─────────
  1    │ Inception.2010.1080p.mkv       │ Inception                      │  2010 │ movie  │   8.4  │ ✅ new
  2    │ Breaking.Bad.S01E01.mkv        │ Breaking Bad                   │  2008 │ tv     │   9.5  │ ✅ new
  3    │ The.Dark.Knight.REMUX.mkv      │ The Dark Knight                │  2008 │ movie  │   9.0  │ ⏩ skip
  ─────┴────────────────────────────────┴────────────────────────────────┴───────┴────────┴────────┴─────────
```

### Column Definitions

| Column       | Width | Align | Description                         |
|--------------|-------|-------|-------------------------------------|
| `#`          | 4     | left  | 1-based index                       |
| `File Name`  | 30    | left  | Original filename (truncated + `…`) |
| `Clean Title`| 30    | left  | Parsed clean title (truncated + `…`)|
| `Year`       | 5     | right | Year or `  -  ` if unknown          |
| `Type`       | 6     | left  | `movie` or `tv`                     |
| `Rating`     | 6     | right | TMDb rating or `   -  ` if unknown  |
| `Status`     | 8     | left  | `✅ new`, `⏩ skip`, or `❌ err`    |

### Truncation
Strings longer than column width are truncated to `width - 1` chars + `…`.

---

## 10. Error Handling in Output

- Errors during scan are logged via `errlog` to:
  1. `.movie-output/logs/error.txt` (flat file with stack traces)
  2. `error_logs` table in SQLite database
- Errors are NOT printed to the main scan output unless they are
  user-facing (e.g., "folder not found")
- Warnings (TMDb miss, thumbnail fail) print inline with the item

---

## 11. REST API Request Logging

When the REST server is running (via `movie rest` or `movie scan --rest`),
every HTTP request is logged via `errlog.Info`:

```
[REST] GET /api/media → 200 (2.3ms)
[REST] PATCH /api/media/5/watched → 200 (1.1ms)
```

Format: `[REST] {METHOD} {PATH} → {STATUS} ({DURATION})`

Logs go to both the file logger and the `error_logs` DB table at `INFO` level.

---

## 12. Exit Codes

| Code | Meaning                                                  |
|------|----------------------------------------------------------|
| `0`  | Success — scan completed (even if some items were skipped or had warnings) |
| `1`  | Fatal error — folder not found, database open failed, or cobra command error |

### Rules
- Individual item failures (TMDb miss, thumbnail download fail, DB insert error) do NOT cause a non-zero exit.
- These are logged via `errlog` and the scan continues.
- Only top-level failures (invalid folder, DB connection error) cause exit 1.
- When `--rest` is used, the process blocks on the HTTP server and exits 0 on Ctrl+C (SIGINT).
