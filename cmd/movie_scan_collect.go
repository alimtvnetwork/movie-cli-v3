// movie_scan_collect.go — video file discovery for movie scan
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/movie-cli-v3/cleaner"
)

// videoFile holds a discovered video file's display name and full path.
type videoFile struct {
	Name     string // display name used for cleaning (dir name or filename)
	FullPath string // absolute path to the actual video file
}

// collectVideoFiles finds video files in the given directory.
// When recursive is true, it walks subdirectories up to maxDepth levels (0 = unlimited).
func collectVideoFiles(scanDir string, recursive bool, maxDepth int) []videoFile {
	var files []videoFile

	if recursive {
		scanDir = filepath.Clean(scanDir)
		baseParts := len(splitPath(scanDir))

		_ = filepath.WalkDir(scanDir, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				fmt.Fprintf(os.Stderr, "  ⚠️  Cannot access %s: %v\n", path, err)
				return nil // continue walking
			}
			// Skip .movie-output and hidden directories
			if d.IsDir() {
				base := d.Name()
				if base == ".movie-output" || (strings.HasPrefix(base, ".") && base != ".") {
					return filepath.SkipDir
				}
				// Enforce depth limit
				if maxDepth > 0 {
					dirParts := len(splitPath(filepath.Clean(path)))
					if dirParts-baseParts > maxDepth {
						return filepath.SkipDir
					}
				}
				return nil
			}
			// Check depth for files too
			if maxDepth > 0 {
				fileParts := len(splitPath(filepath.Clean(filepath.Dir(path))))
				if fileParts-baseParts > maxDepth {
					return nil
				}
			}
			if cleaner.IsVideoFile(d.Name()) {
				// Use parent directory name if it differs from scanDir, else use filename
				parentDir := filepath.Dir(path)
				name := d.Name()
				if parentDir != scanDir {
					name = filepath.Base(parentDir)
				}
				files = append(files, videoFile{Name: name, FullPath: path})
			}
			return nil
		})
	} else {
		entries, readErr := os.ReadDir(scanDir)
		if readErr != nil {
			fmt.Fprintf(os.Stderr, "❌ Cannot read folder: %v\n", readErr)
			return nil
		}
		for _, entry := range entries {
			name := entry.Name()
			fullPath := filepath.Join(scanDir, name)

			if entry.IsDir() {
				// Look for first video file inside the subdirectory
				subEntries, subErr := os.ReadDir(fullPath)
				if subErr != nil {
					fmt.Fprintf(os.Stderr, "  ⚠️  Cannot read subdirectory %s: %v\n", name, subErr)
					continue
				}
				for _, sub := range subEntries {
					if !sub.IsDir() && cleaner.IsVideoFile(sub.Name()) {
						files = append(files, videoFile{
							Name:     entry.Name(),
							FullPath: filepath.Join(fullPath, sub.Name()),
						})
						break
					}
				}
			} else if cleaner.IsVideoFile(name) {
				files = append(files, videoFile{Name: name, FullPath: fullPath})
			}
		}
	}

	return files
}

// splitPath splits a filepath into its components.
func splitPath(p string) []string {
	var parts []string
	for p != "" && p != "." && p != "/" && p != string(filepath.Separator) {
		dir, file := filepath.Split(p)
		if file != "" {
			parts = append(parts, file)
		}
		p = filepath.Clean(dir)
		if p == dir {
			break
		}
	}
	return parts
}
