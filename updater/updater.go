package updater

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Result holds the outcome of a self-update attempt.
type Result struct {
	PreviousVersion string
	UpdatedTo       string
	AfterCommit     string
	RepoPath        string
	Output          string
	AlreadyLatest   bool
}

func Run() (*Result, error) {
	if _, err := exec.LookPath("git"); err != nil {
		return nil, fmt.Errorf("git is not installed or not in PATH")
	}

	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("cannot read current directory: %w", err)
	}

	repoPath, err := gitOutput(cwd, "rev-parse", "--show-toplevel")
	if err != nil {
		return nil, fmt.Errorf("self-update must run inside a cloned git repository: %w", err)
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
