// Package templates provides embedded HTML templates for the movie CLI.
package templates

import "embed"

//go:embed report.html
var FS embed.FS
