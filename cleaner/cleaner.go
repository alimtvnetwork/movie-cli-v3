// Package cleaner handles filename cleaning and metadata extraction from file names.
package cleaner

import (
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// video file extensions
var videoExtensions = map[string]bool{
	".mkv": true, ".mp4": true, ".avi": true, ".mov": true,
	".wmv": true, ".flv": true, ".webm": true, ".m4v": true,
	".ts": true, ".vob": true, ".ogv": true, ".mpg": true,
	".mpeg": true, ".3gp": true,
}

// patterns to remove from file names
var junkPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)\b(1080p|720p|480p|2160p|4k)\b`),
	regexp.MustCompile(`(?i)\b(bluray|bdrip|brrip|webrip|web-dl|web\.dl|hdrip|dvdrip|dvdscr|hdtv|hdcam|cam|ts|tc)\b`),
	regexp.MustCompile(`(?i)\b(x264|x265|h\.?264|h\.?265|hevc|aac|ac3|dts|mp3|flac|atmos)\b`),
	regexp.MustCompile(`(?i)\b(rarbg|yts|yify|eztv|ettv|sparks|geckos|fgt|evo|ion10|cmrg|ntb|megusta)\b`),
	regexp.MustCompile(`(?i)\b(extended|unrated|directors\.?cut|remastered|proper|rerip|internal)\b`),
	regexp.MustCompile(`(?i)\b(multi|dual|eng|english|hindi|spanish|french|german|ita|sub|subs|subtitle|subtitles)\b`),
	regexp.MustCompile(`\[.*?\]`),
	regexp.MustCompile(`\(.*?\)`),
	regexp.MustCompile(`\{.*?\}`),
}

var yearPattern = regexp.MustCompile(`\b((?:19|20)\d{2})\b`)
var tvPattern = regexp.MustCompile(`(?i)S\d{1,2}E\d{1,2}|Season\s*\d+|Episode\s*\d+`)
var multiSpace = regexp.MustCompile(`\s{2,}`)

// Result holds the cleaned information extracted from a filename.
type Result struct {
	CleanTitle string
	Extension  string
	Type       string // "movie" or "tv"
	Year       int
}

// IsVideoFile checks if the file has a video extension.
func IsVideoFile(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	return videoExtensions[ext]
}

// Clean takes a raw filename and returns cleaned metadata.
func Clean(filename string) Result {
	ext := filepath.Ext(filename)
	name := strings.TrimSuffix(filename, ext)

	// Detect TV show
	mediaType := "movie"
	if tvPattern.MatchString(name) {
		mediaType = "tv"
	}

	// Extract year
	year := 0
	if matches := yearPattern.FindStringSubmatch(name); len(matches) > 1 {
		year, _ = strconv.Atoi(matches[1])
	}

	// Replace dots and underscores with spaces
	cleaned := strings.ReplaceAll(name, ".", " ")
	cleaned = strings.ReplaceAll(cleaned, "_", " ")

	// Remove junk patterns
	for _, p := range junkPatterns {
		cleaned = p.ReplaceAllString(cleaned, "")
	}

	// Clean up extra spaces
	cleaned = multiSpace.ReplaceAllString(cleaned, " ")
	cleaned = strings.TrimSpace(cleaned)

	// If year exists, truncate everything after the year
	if year > 0 {
		yearStr := strconv.Itoa(year)
		if idx := strings.Index(cleaned, yearStr); idx > 0 {
			cleaned = strings.TrimSpace(cleaned[:idx+len(yearStr)])
		}
	}

	return Result{
		CleanTitle: cleaned,
		Year:       year,
		Extension:  strings.ToLower(ext),
		Type:       mediaType,
	}
}

// ToSlug converts a clean title to a slug (e.g., "Scream 2022" → "scream-2022").
func ToSlug(title string) string {
	s := strings.ToLower(title)
	s = regexp.MustCompile(`[^a-z0-9\s-]`).ReplaceAllString(s, "")
	s = multiSpace.ReplaceAllString(s, " ")
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, " ", "-")
	return s
}

// ToCleanFileName creates a clean filename: "Scream (2022).mkv"
func ToCleanFileName(title string, year int, ext string) string {
	if year > 0 {
		return title + " (" + strconv.Itoa(year) + ")" + ext
	}
	return title + ext
}
