// Package version holds the three build-time variables that describe exactly
// which binary is running.
//
// The values are injected at compile time with -ldflags so every released
// binary carries its own version stamp:
//
//	go build -ldflags "
//	  -X github.com/alimtvnetwork/movie-cli-v3/version.Version=v1.2.0
//	  -X github.com/alimtvnetwork/movie-cli-v3/version.Commit=abc1234
//	  -X github.com/alimtvnetwork/movie-cli-v3/version.BuildDate=2024-06-01
//	" .
//
// During development (no -ldflags) the defaults below are used instead.
package version

import "fmt"

// These three variables are overwritten by -ldflags at build time.
var (
	Version   = "v1.17.0" // semver tag  e.g. "v1.2.0"
	Commit    = "none"       // git SHA     e.g. "abc1234"
	BuildDate = "unknown"    // build date  e.g. "2024-06-01"
)

// Full returns the full version string printed by `movie version`.
//
//	v1.2.0 (commit: abc1234, built: 2024-06-01)
func Full() string {
	return fmt.Sprintf("%s (commit: %s, built: %s)", Version, Commit, BuildDate)
}

// Short returns just the semver tag, e.g. "v1.2.0".
// Used by the updater when comparing against the latest GitHub release tag.
func Short() string {
	return Version
}
