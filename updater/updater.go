package updater

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// repoURL is the canonical GitHub URL used when the binary is not inside the repo.
const repoURL = "https://github.com/alimtvnetwork/movie-cli-v3.git"

// Result holds the outcome of a self-update attempt.
type Result struct {
	PreviousVersion string
	UpdatedTo       string
	AfterCommit     string
	RepoPath        string
	Output          string
	AlreadyLatest   bool
	Bootstrapped    bool
}

func Run() (*Result, error) {
	if _, err := exec.LookPath("git"); err != nil {
		return nil, fmt.Errorf("git is not installed or not in PATH")
	}

	repoPath, bootstrapped, err := findRepoPath()
	if err != nil {
		return nil, err
	}

	if bootstrapped {
		afterCommit, err := gitOutput(repoPath, "rev-parse", "--short", "HEAD")
		if err != nil {
			return nil, fmt.Errorf("cannot read bootstrapped commit: %w", err)
		}
		return &Result{
			Bootstrapped: true,
			UpdatedTo:    afterCommit,
			AfterCommit:  afterCommit,
			RepoPath:     repoPath,
		}, nil
	}

	dirty, err := gitOutput(repoPath, "status", "--porcelain")
	if err != nil {
		return nil, fmt.Errorf("cannot check git status: %w", err)
	}
	if strings.TrimSpace(dirty) != "" {
		return nil, fmt.Errorf("repository has local changes; commit or stash them before self-update")
	}

	beforeCommit, err := gitOutput(repoPath, "rev-parse", "--short", "HEAD")
	if err != nil {
		return nil, fmt.Errorf("cannot read current commit: %w", err)
	}

	pullOutput, err := gitOutput(repoPath, "pull", "--ff-only")
	if err != nil {
		return nil, fmt.Errorf("git pull failed: %w", err)
	}

	afterCommit, err := gitOutput(repoPath, "rev-parse", "--short", "HEAD")
	if err != nil {
		return nil, fmt.Errorf("cannot read updated commit: %w", err)
	}

	updated := beforeCommit != afterCommit

	return &Result{
		AlreadyLatest:   !updated,
		PreviousVersion: beforeCommit,
		UpdatedTo:       afterCommit,
		AfterCommit:     afterCommit,
		RepoPath:        repoPath,
		Output:          pullOutput,
	}, nil
}

func gitOutput(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir

	out, err := cmd.CombinedOutput()
	text := strings.TrimSpace(string(out))
	if err != nil {
		if text == "" {
			return "", err
		}
		return "", fmt.Errorf("%s", text)
	}

	return text, nil
}

// findRepoPath locates the git repository root by checking (in order):
//  1. The directory containing the running binary
//  2. The current working directory
//  3. A default clone location next to the binary
//
// If none have a .git directory, it clones the repo fresh next to the binary.
func findRepoPath() (string, bool, error) {
	// 1. Try the binary's own directory
	exe, exeErr := os.Executable()
	if exeErr == nil {
		exe, _ = filepath.EvalSymlinks(exe) // resolve symlinks
		exeDir := filepath.Dir(exe)
		if p, gitErr := gitOutput(exeDir, "rev-parse", "--show-toplevel"); gitErr == nil {
			return p, false, nil
		}

		// Check for a clone next to the binary: <exeDir>/movie-cli-v3/
		cloneDir := filepath.Join(exeDir, "movie-cli-v3")
		if p, gitErr := gitOutput(cloneDir, "rev-parse", "--show-toplevel"); gitErr == nil {
			return p, false, nil
		}
	}

	// 2. Try CWD
	cwd, cwdErr := os.Getwd()
	if cwdErr == nil {
		if p, gitErr := gitOutput(cwd, "rev-parse", "--show-toplevel"); gitErr == nil {
			return p, false, nil
		}
	}

	// 3. No repo found — clone next to the binary
	if exeErr == nil {
		exeDir := filepath.Dir(exe)
		cloneDir := filepath.Join(exeDir, "movie-cli-v3")
		fmt.Printf("📥 No local repo found. Cloning to: %s\n", cloneDir)
		if _, cloneErr := gitOutput(exeDir, "clone", "--depth", "1", repoURL); cloneErr != nil {
			return "", false, fmt.Errorf("cannot clone repository: %w", cloneErr)
		}
		return cloneDir, true, nil
	}

	return "", false, fmt.Errorf("cannot locate the movie-cli-v3 repository. Run from the repo directory or ensure the binary is deployed alongside the repo")
}
