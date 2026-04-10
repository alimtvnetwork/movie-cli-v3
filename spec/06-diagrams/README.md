# Diagrams Index

All architectural and flow diagrams for the `movie` CLI.

| # | File | Description |
|---|------|-------------|
| 01 | [scan-flow](images/01-scan-flow.png) | How `movie scan` discovers files, fetches TMDb metadata, and writes to DB |
| 02 | [ls-flow](images/02-ls-flow.png) | How `movie ls` queries the media table with filters and renders output |
| 03 | [search-flow](images/03-search-flow.png) | How `movie search` performs fuzzy/exact lookups across local DB and TMDb |
| 04 | [info-flow](images/04-info-flow.png) | How `movie info` resolves a movie and displays detailed metadata |
| 05 | [suggest-flow](images/05-suggest-flow.png) | How `movie suggest` picks random or filtered recommendations |
| 06 | [move-flow](images/06-move-flow.png) | How `movie move` relocates a file and records history |
| 07 | [rename-flow](images/07-rename-flow.png) | How `movie rename` batch-renames files to match TMDb titles |
| 08 | [undo-flow](images/08-undo-flow.png) | How `movie undo` reverses the last move/rename operation |
| 09 | [play-flow](images/09-play-flow.png) | How `movie play` resolves and launches a media file |
| 10 | [stats-flow](images/10-stats-flow.png) | How `movie stats` aggregates and displays collection statistics |
| 11 | [config-flow](images/11-config-flow.png) | How `movie config` reads/writes settings in the config table |
| 12 | [tag-flow](images/12-tag-flow.png) | How `movie tag` manages custom labels via the tags/media_tags tables |
| 13 | [tmdb-data-flow](images/13-tmdb-data-flow.png) | End-to-end TMDb API interaction: search → fetch → cache → store |
| 14 | [cleaner-flow](images/14-cleaner-flow.png) | How `movie cleaner` detects and removes orphaned/duplicate DB entries |
| 15 | [db-schema](images/15-db-schema.png) | Full database schema with all tables and relationships |
| 16 | [command-db-access](images/16-command-db-access.png) | Matrix showing which commands read/write which tables |
| 17 | [move-undo-lifecycle](images/17-move-undo-lifecycle.png) | Move → undo lifecycle with history state transitions |
| 18 | [rename-vs-move-lifecycle](images/18-rename-vs-move-lifecycle.png) | Rename vs move comparison — shared history tracking, DB read/write patterns |

## Source Files

Mermaid source files (`.mmd`) are stored alongside this README for regeneration.
