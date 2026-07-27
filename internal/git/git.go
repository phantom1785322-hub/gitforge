// Package git provides Git operations using the Git CLI.
package git

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Repo represents a Git repository.
type Repo struct {
	path    string
	workdir string
}

// OpenRepo opens an existing repository at the given path.
func OpenRepo(path string) (*Repo, error) {
	_, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("absolute path: %w", err)
	}

	gitDir := filepath.Join(path, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		if !isBareRepo(path) {
			return nil, fmt.Errorf("not a git repository: %s", path)
		}
	}

	workdir := path
	if !isBareRepo(path) {
		workdir, _ = filepath.Abs(path)
	}

	return &Repo{
		path:    path,
		workdir: workdir,
	}, nil
}

// InitRepo initializes a new repository at the given path.
func InitRepo(path string, bare bool) (*Repo, error) {
	_, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("absolute path: %w", err)
	}

	args := []string{"init"}
	if bare {
		args = append(args, "--bare")
	}
	args = append(args, path)

	cmd := exec.Command("git", args...)
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("init repository: %w", err)
	}

	repo, err := OpenRepo(path)
	if err != nil {
		return nil, err
	}

	repo.SetConfig("user.name", "GitForge User")
	repo.SetConfig("user.email", "user@gitforge.local")
	repo.SetConfig("init.defaultBranch", "main")
	repo.SetConfig("push.autoSetupRemote", "true")
	repo.SetConfig("core.autocrlf", "false")
	repo.SetConfig("core.filemode", "true")

	return repo, nil
}

// CloneOptions configures repository cloning.
type CloneOptions struct {
	Branch            string
	Depth             int
	RecurseSubmodules bool
	SingleBranch      bool
}

// CloneRepo clones a repository from URL to path.
func CloneRepo(ctx context.Context, url, path string, opts *CloneOptions) (*Repo, error) {
	_, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("absolute path: %w", err)
	}

	args := []string{"clone"}

	if opts != nil {
		if opts.Branch != "" {
			args = append(args, "--branch", opts.Branch)
		}
		if opts.Depth > 0 {
			args = append(args, "--depth", strconv.Itoa(opts.Depth))
		}
		if opts.RecurseSubmodules {
			args = append(args, "--recurse-submodules")
		}
		if opts.SingleBranch {
			args = append(args, "--single-branch")
		}
	}

	args = append(args, url, path)

	cmd := exec.CommandContext(ctx, "git", args...)
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("clone repository: %w", err)
	}

	return OpenRepo(path)
}

// isBareRepo checks if a directory is a bare git repository.
func isBareRepo(path string) bool {
	required := []string{"HEAD", "config", "objects", "refs"}
	for _, f := range required {
		if _, err := os.Stat(filepath.Join(path, f)); os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// isWorktree checks if .git is a worktree pointer file.
func isWorktree(gitPath string) bool {
	data, err := os.ReadFile(gitPath)
	if err != nil {
		return false
	}
	return strings.HasPrefix(string(data), "gitdir:")
}

// OpenRepoFromPath opens an existing repository.
func OpenRepoFromPath(path string) (*Repo, error) {
	_, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("absolute path: %w", err)
	}

	gitDir := filepath.Join(path, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		if !isBareRepo(path) {
			return nil, fmt.Errorf("not a git repository: %s", path)
		}
	} else {
		if isWorktree(gitDir) {
			// Handle worktree case
		}
	}

	return &Repo{
		path:    path,
		workdir: path,
	}, nil
}

func (r *Repo) Path() string {
	return r.path
}

// Workdir returns the working directory.
func (r *Repo) Workdir() string {
	return r.workdir
}

// IsBare returns true if the repository is bare.
func (r *Repo) IsBare() bool {
	return isBareRepo(r.path)
}

// git runs a git command and returns stdout.
func (r *Repo) git(args ...string) (string, error) {
	return r.gitWithContext(context.Background(), args...)
}

// gitWithContext runs a git command with context.
func (r *Repo) gitWithContext(ctx context.Context, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = r.path
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git %s: %w (stderr: %s)", strings.Join(args, " "), err, stderr.String())
	}

	return strings.TrimSpace(stdout.String()), nil
}

// Head returns the current HEAD reference.
func (r *Repo) Head() (string, error) {
	return r.git("rev-parse", "--abbrev-ref", "HEAD")
}

// Branch returns the current branch name.
func (r *Repo) Branch() (string, error) {
	return r.Head()
}

// Branches returns all local branches.
func (r *Repo) Branches() ([]Branch, error) {
	output, err := r.git("branch", "--format=%(refname:short)%(if:equals=*) %(upstream:track)%(then)%(upstream:short)%(end)")
	if err != nil {
		return nil, err
	}

	var branches []Branch
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, " ", 2)
		branch := Branch{
			Name: parts[0],
		}
		if len(parts) > 1 {
			branch.Upstream = parts[1]
			branch.Tracking = true
		}
		branches = append(branches, branch)
	}
	return branches, nil
}

// Branch represents a git branch.
type Branch struct {
	Name      string
	Upstream  string
	Tracking  bool
	IsCurrent bool
}

// RemoteBranches returns all remote branches.
func (r *Repo) RemoteBranches() ([]Branch, error) {
	output, err := r.git("branch", "-r", "--format=%(refname:short)")
	if err != nil {
		return nil, err
	}

	var branches []Branch
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		branches = append(branches, Branch{
			Name:     line,
			Tracking: true,
		})
	}
	return branches, nil
}

// CreateBranch creates a new branch.
func (r *Repo) CreateBranch(name, startPoint string) error {
	args := []string{"branch"}
	if startPoint != "" {
		args = append(args, startPoint)
	}
	args = append(args, name)
	_, err := r.git(args...)
	return err
}

// DeleteBranch deletes a branch.
func (r *Repo) DeleteBranch(name string, force bool) error {
	args := []string{"branch"}
	if force {
		args = append(args, "-D")
	} else {
		args = append(args, "-d")
	}
	args = append(args, name)
	_, err := r.git(args...)
	return err
}

// Checkout checks out a branch or commit.
func (r *Repo) Checkout(target string) error {
	_, err := r.git("checkout", target)
	return err
}

// CheckoutNewBranch creates and checks out a new branch.
func (r *Repo) CheckoutNewBranch(name, startPoint string) error {
	args := []string{"checkout", "-b", name}
	if startPoint != "" {
		args = append(args, startPoint)
	}
	_, err := r.git(args...)
	return err
}

// Status returns the repository status.
func (r *Repo) Status() (Status, error) {
	output, err := r.git("status", "--porcelain=v2")
	if err != nil {
		return Status{}, err
	}

	var status Status
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "1 ") {
			parts := strings.SplitN(line, " ", 9)
			if len(parts) >= 9 {
				file := parts[8]
				switch parts[1] {
				case "A", "M", "D", "R", "C":
					status.Staged = append(status.Staged, file)
				}
			}
		} else if strings.HasPrefix(line, "2 ") {
			parts := strings.SplitN(line, " ", 10)
			if len(parts) >= 10 {
				file := parts[9]
				status.Staged = append(status.Staged, file)
			}
		} else if strings.HasPrefix(line, "?") {
			file := strings.TrimPrefix(line, "? ")
			status.Untracked = append(status.Untracked, file)
		} else if strings.HasPrefix(line, "u ") {
			parts := strings.SplitN(line, " ", 3)
			if len(parts) >= 3 {
				status.Conflicted = append(status.Conflicted, parts[2])
			}
		} else if strings.HasPrefix(line, "1 ") && strings.Contains(line, " ") {
			parts := strings.SplitN(line, " ", 9)
			if len(parts) >= 9 {
				file := parts[8]
				if parts[2] != "." {
					status.Unstaged = append(status.Unstaged, file)
				}
			}
		}
	}

	return status, nil
}

// Status represents repository status.
type Status struct {
	Staged      []string
	Unstaged    []string
	Untracked   []string
	Conflicted  []string
}

// IsClean returns true if the working directory is clean.
func (s Status) IsClean() bool {
	return len(s.Staged) == 0 && len(s.Unstaged) == 0 && len(s.Untracked) == 0 && len(s.Conflicted) == 0
}

// Add adds files to the index.
func (r *Repo) Add(paths ...string) error {
	args := append([]string{"add"}, paths...)
	_, err := r.git(args...)
	return err
}

// Remove removes files from the index.
func (r *Repo) Remove(paths ...string) error {
	args := append([]string{"rm", "--cached"}, paths...)
	_, err := r.git(args...)
	return err
}

// Commit creates a commit with the given message.
func (r *Repo) Commit(message string, author *Signature) (string, error) {
	args := []string{"commit", "-m", message}

	if author != nil {
		args = append(args, "--author", fmt.Sprintf("%s <%s>", author.Name, author.Email))
	}

	output, err := r.git(args...)
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`\[([a-f0-9]+)\]`)
	matches := re.FindStringSubmatch(output)
	if len(matches) > 1 {
		return matches[1], nil
	}

	hash, err := r.git("rev-parse", "HEAD")
	if err != nil {
		return "", err
	}
	return hash, nil
}

// Signature represents a commit author/committer.
type Signature struct {
	Name  string
	Email string
	When  time.Time
}

// Log returns commit history.
func (r *Repo) Log(opts *LogOptions) ([]Commit, error) {
	if opts == nil {
		opts = &LogOptions{}
	}

	args := []string{"log", "--pretty=format:%H|%an|%ae|%at|%s|%b", "--date=unix"}

	if opts.Limit > 0 {
		args = append(args, fmt.Sprintf("-%d", opts.Limit))
	}
	if opts.Since != "" {
		args = append(args, "--since", opts.Since)
	}
	if opts.Until != "" {
		args = append(args, "--until", opts.Until)
	}
	if opts.Author != "" {
		args = append(args, "--author", opts.Author)
	}
	if opts.Grep != "" {
		args = append(args, "--grep", opts.Grep)
	}
	if opts.From != "" {
		args = append(args, opts.From+"..HEAD")
	}

	output, err := r.git(args...)
	if err != nil {
		return nil, err
	}

	var commits []Commit
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "|", 6)
		if len(parts) < 5 {
			continue
		}

		timestamp, _ := strconv.ParseInt(parts[3], 10, 64)
		commit := Commit{
			Hash:     parts[0],
			Author:   Signature{Name: parts[1], Email: parts[2]},
			Date:     time.Unix(timestamp, 0),
			Subject:  parts[4],
			Body:     parts[5],
		}
		commits = append(commits, commit)
	}

	return commits, nil
}

// LogOptions configures log output.
type LogOptions struct {
	Limit  int
	Since  string
	Until  string
	Author string
	Grep   string
	From   string
}

// Commit represents a git commit.
type Commit struct {
	Hash    string
	Author  Signature
	Date    time.Time
	Subject string
	Body    string
}

// Diff returns the diff between two commits or working tree.
func (r *Repo) Diff(opts *DiffOptions) (*Diff, error) {
	if opts == nil {
		opts = &DiffOptions{}
	}

	args := []string{"diff"}

	if opts.Cached {
		args = append(args, "--cached")
	}
	if opts.From != "" {
		args = append(args, opts.From)
	}
	if opts.To != "" {
		args = append(args, opts.To)
	}
	if opts.File != "" {
		args = append(args, "--", opts.File)
	}

	output, err := r.git(args...)
	if err != nil {
		return nil, err
	}

	diff := &Diff{Raw: output}
	diff.parse()
	return diff, nil
}

// DiffOptions configures diff output.
type DiffOptions struct {
	From   string
	To     string
	Cached bool
	File   string
}

// Diff represents a git diff.
type Diff struct {
	Raw     string
	Patches []Patch
}

// Patch represents a single file diff.
type Patch struct {
	OldFile string
	NewFile string
	Status  string
	Hunks   []Hunk
}

// Hunk represents a diff hunk.
type Hunk struct {
	OldStart int
	OldLines int
	NewStart int
	NewLines int
	Lines    []DiffLine
}

// DiffLine represents a line in a diff.
type DiffLine struct {
	Type    string // "add", "del", "ctx"
	Content string
	OldLine int
	NewLine int
}

// parse parses the raw diff into structured patches.
func (d *Diff) parse() {
	lines := strings.Split(d.Raw, "\n")
	var currentPatch *Patch
	var currentHunk *Hunk

	for _, line := range lines {
		if strings.HasPrefix(line, "diff --git ") {
			if currentPatch != nil {
				d.Patches = append(d.Patches, *currentPatch)
			}
			currentPatch = &Patch{}
			parts := strings.Fields(line)
			if len(parts) >= 4 {
				currentPatch.OldFile = strings.TrimPrefix(parts[2], "a/")
				currentPatch.NewFile = strings.TrimPrefix(parts[3], "b/")
				currentPatch.Status = "modified"
			}
		} else if strings.HasPrefix(line, "new file mode") {
			if currentPatch != nil {
				currentPatch.Status = "added"
			}
		} else if strings.HasPrefix(line, "deleted file mode") {
			if currentPatch != nil {
				currentPatch.Status = "deleted"
			}
		} else if strings.HasPrefix(line, "rename from") {
			if currentPatch != nil {
				currentPatch.Status = "renamed"
			}
		} else if strings.HasPrefix(line, "@@ ") {
			if currentPatch != nil {
				if currentHunk != nil {
					currentPatch.Hunks = append(currentPatch.Hunks, *currentHunk)
				}
				currentHunk = &Hunk{}
				re := regexp.MustCompile(`@@ -(\d+),?(\d*) \+(\d+),?(\d*) @@`)
				matches := re.FindStringSubmatch(line)
				if len(matches) >= 5 {
					currentHunk.OldStart, _ = strconv.Atoi(matches[1])
					if matches[2] != "" {
						currentHunk.OldLines, _ = strconv.Atoi(matches[2])
					} else {
						currentHunk.OldLines = 1
					}
					currentHunk.NewStart, _ = strconv.Atoi(matches[3])
					if matches[4] != "" {
						currentHunk.NewLines, _ = strconv.Atoi(matches[4])
					} else {
						currentHunk.NewLines = 1
					}
				}
			}
		} else if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
			if currentHunk != nil {
				currentHunk.Lines = append(currentHunk.Lines, DiffLine{
					Type:    "add",
					Content: line[1:],
				})
			}
		} else if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
			if currentHunk != nil {
				currentHunk.Lines = append(currentHunk.Lines, DiffLine{
					Type:    "del",
					Content: line[1:],
				})
			}
		} else if strings.HasPrefix(line, " ") && currentHunk != nil {
			currentHunk.Lines = append(currentHunk.Lines, DiffLine{
				Type:    "ctx",
				Content: line[1:],
			})
		}
	}

	if currentPatch != nil {
		if currentHunk != nil {
			currentPatch.Hunks = append(currentPatch.Hunks, *currentHunk)
		}
		d.Patches = append(d.Patches, *currentPatch)
	}
}

// Remotes returns all configured remotes.
func (r *Repo) Remotes() ([]Remote, error) {
	output, err := r.git("remote", "-v")
	if err != nil {
		return nil, err
	}

	seen := make(map[string]bool)
	var remotes []Remote
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) >= 3 {
			name := parts[0]
			url := parts[1]
			direction := parts[2]

			if !seen[name] {
				remote := Remote{Name: name}
				seen[name] = true
				if strings.HasSuffix(direction, "(fetch)") {
					remote.FetchURL = strings.TrimSuffix(url, " (fetch)")
				} else if strings.HasSuffix(direction, "(push)") {
					remote.PushURL = strings.TrimSuffix(url, " (push)")
				}
				remotes = append(remotes, remote)
			} else {
				for i := range remotes {
					if remotes[i].Name == name {
						if strings.HasSuffix(direction, "(fetch)") {
							remotes[i].FetchURL = strings.TrimSuffix(url, " (fetch)")
						} else if strings.HasSuffix(direction, "(push)") {
							remotes[i].PushURL = strings.TrimSuffix(url, " (push)")
						}
						break
					}
				}
			}
		}
	}
	return remotes, nil
}

// Remote represents a git remote.
type Remote struct {
	Name     string
	FetchURL string
	PushURL  string
}

// AddRemote adds a new remote.
func (r *Repo) AddRemote(name, url string) error {
	_, err := r.git("remote", "add", name, url)
	return err
}

// RemoveRemote removes a remote.
func (r *Repo) RemoveRemote(name string) error {
	_, err := r.git("remote", "remove", name)
	return err
}

// Fetch fetches from a remote.
func (r *Repo) Fetch(remote string, refspecs []string) error {
	args := []string{"fetch", remote}
	args = append(args, refspecs...)
	_, err := r.git(args...)
	return err
}

// Push pushes to a remote.
func (r *Repo) Push(remote string, refspecs []string) error {
	args := []string{"push", remote}
	args = append(args, refspecs...)
	_, err := r.git(args...)
	return err
}

// Pull pulls from a remote.
func (r *Repo) Pull(remote, branch string) error {
	args := []string{"pull", remote}
	if branch != "" {
		args = append(args, branch)
	}
	_, err := r.git(args...)
	return err
}

// Merge merges a branch into the current branch.
func (r *Repo) Merge(branch string) error {
	_, err := r.git("merge", branch)
	return err
}

// Reset resets the current branch to a commit.
func (r *Repo) Reset(target string, mode string) error {
	args := []string{"reset"}
	if mode != "" {
		args = append(args, "--"+mode)
	}
	args = append(args, target)
	_, err := r.git(args...)
	return err
}

// Stash saves changes to a stash.
func (r *Repo) Stash(message string) error {
	args := []string{"stash", "push"}
	if message != "" {
		args = append(args, "-m", message)
	}
	_, err := r.git(args...)
	return err
}

// StashList returns all stashes.
func (r *Repo) StashList() ([]Stash, error) {
	output, err := r.git("stash", "list", "--format=%gd|%gs|%gD")
	if err != nil {
		return nil, err
	}

	var stashes []Stash
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Split(line, "|")
		if len(parts) >= 3 {
			stashes = append(stashes, Stash{
				Index:   parts[0],
				Message: parts[1],
				Date:    parts[2],
			})
		}
	}
	return stashes, nil
}

// Stash represents a git stash.
type Stash struct {
	Index   string
	Message string
	Date    string
}

// StashPop pops the latest stash.
func (r *Repo) StashPop() error {
	_, err := r.git("stash", "pop")
	return err
}

// StashDrop drops a stash.
func (r *Repo) StashDrop(index int) error {
	_, err := r.git("stash", "drop", fmt.Sprintf("stash@{%d}", index))
	return err
}

// Config returns the repository config.
func (r *Repo) Config() (*Config, error) {
	output, err := r.git("config", "--list", "--show-origin")
	if err != nil {
		return nil, err
	}

	cfg := &Config{Values: make(map[string]string)}
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "\t") {
			parts := strings.SplitN(line, "\t", 2)
			if len(parts) == 2 {
				cfg.Values[parts[0]] = parts[1]
			}
		}
	}
	return cfg, nil
}

// Config represents git configuration.
type Config struct {
	Values map[string]string
}

// GetConfig gets a config value.
func (r *Repo) GetConfig(key string) (string, error) {
	return r.git("config", "--get", key)
}

// SetConfig sets a config value.
func (r *Repo) SetConfig(key, value string) error {
	_, err := r.git("config", key, value)
	return err
}

// Submodules returns all submodules.
func (r *Repo) Submodules() ([]Submodule, error) {
	output, err := r.git("submodule", "status")
	if err != nil {
		return nil, err
	}

	var submodules []Submodule
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			submodules = append(submodules, Submodule{
				Commit: parts[0],
				Path:   parts[1],
			})
		}
	}
	return submodules, nil
}

// Submodule represents a git submodule.
type Submodule struct {
	Commit string
	Path   string
	Name   string
	URL    string
}

// AddSubmodule adds a submodule.
func (r *Repo) AddSubmodule(url, path string) error {
	_, err := r.git("submodule", "add", url, path)
	return err
}

// Tags returns all tags.
func (r *Repo) Tags() ([]Tag, error) {
	output, err := r.git("tag", "-l", "--format=%(refname:short)|%(creatordate:short)|%(subject)")
	if err != nil {
		return nil, err
	}

	var tags []Tag
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Split(line, "|")
		if len(parts) >= 3 {
			tags = append(tags, Tag{
				Name:    parts[0],
				Date:    parts[1],
				Message: parts[2],
			})
		}
	}
	return tags, nil
}

// Tag represents a git tag.
type Tag struct {
	Name    string
	Date    string
	Message string
}

// CreateTag creates a new tag.
func (r *Repo) CreateTag(name, commit, message string) error {
	args := []string{"tag", "-a", name, "-m", message}
	if commit != "" {
		args = append(args, commit)
	}
	_, err := r.git(args...)
	return err
}

// DeleteTag deletes a tag.
func (r *Repo) DeleteTag(name string) error {
	_, err := r.git("tag", "-d", name)
	return err
}

// Worktree returns the worktree path.
func (r *Repo) Worktree() string {
	return r.workdir
}

// Close closes the repository (no-op for CLI wrapper).
func (r *Repo) Close() {}

// FindRepoRoot finds the repository root from a path.
func FindRepoRoot(path string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	cmd.Dir = path
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// IsRepo checks if a path is a git repository.
func IsRepo(path string) bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	cmd.Dir = path
	return cmd.Run() == nil
}

// Init initializes a new repository with default config.
func Init(path string) error {
	cmd := exec.Command("git", "init", path)
	if err := cmd.Run(); err != nil {
		return err
	}

	repo, err := OpenRepo(path)
	if err != nil {
		return err
	}

	repo.SetConfig("user.name", "GitForge User")
	repo.SetConfig("user.email", "user@gitforge.local")
	repo.SetConfig("init.defaultBranch", "main")
	repo.SetConfig("push.autoSetupRemote", "true")
	repo.SetConfig("core.autocrlf", "false")
	repo.SetConfig("core.filemode", "true")

	return nil
}

// GetRemoteUrl gets the URL of a remote.
func (r *Repo) GetRemoteUrl(name string) (string, error) {
	return r.git("config", "--get", fmt.Sprintf("remote.%s.url", name))
}

// GetRemotePushUrl gets the push URL of a remote.
func (r *Repo) GetRemotePushUrl(name string) (string, error) {
	return r.git("config", "--get", fmt.Sprintf("remote.%s.pushurl", name))
}

// SetRemoteUrl sets the URL of a remote.
func (r *Repo) SetRemoteUrl(name, url string) error {
	_, err := r.git("remote", "set-url", name, url)
	return err
}

// SetRemotePushUrl sets the push URL of a remote.
func (r *Repo) SetRemotePushUrl(name, url string) error {
	_, err := r.git("remote", "set-url", "--push", name, url)
	return err
}

// GetMergeBase finds the merge base between two commits.
func (r *Repo) GetMergeBase(commit1, commit2 string) (string, error) {
	return r.git("merge-base", commit1, commit2)
}

// CherryPick cherry-picks a commit.
func (r *Repo) CherryPick(commit string) error {
	_, err := r.git("cherry-pick", commit)
	return err
}

// Revert reverts a commit.
func (r *Repo) Revert(commit string) error {
	_, err := r.git("revert", "--no-commit", commit)
	return err
}

// Describe returns a human-readable description of a commit.
func (r *Repo) Describe(commit string) (string, error) {
	return r.git("describe", "--tags", "--always", commit)
}

// Blame returns blame information for a file.
func (r *Repo) Blame(filePath string) ([]BlameHunk, error) {
	output, err := r.git("blame", "-L", "1,$", filePath)
	if err != nil {
		return nil, err
	}

	var hunks []BlameHunk
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		re := regexp.MustCompile(`^(\w+)\s+\(([^)]+)\s+(\d{4}-\d{2}-\d{2}\s+\d{2}:\d{2}:\d{2}\s+\+\d{4})\s+(\d+)\)\s*(.*)`)
		matches := re.FindStringSubmatch(line)
		if len(matches) >= 6 {
			hunks = append(hunks, BlameHunk{
				Commit:  matches[1],
				Author:  matches[2],
				Date:    matches[3],
				LineNum: matches[4],
				Content: matches[5],
			})
		}
	}
	return hunks, nil
}

// BlameHunk represents a blame hunk.
type BlameHunk struct {
	Commit  string
	Author  string
	Date    string
	LineNum string
	Content string
}

// GetFileContent gets the content of a file at a commit.
func (r *Repo) GetFileContent(commit, filePath string) (string, error) {
	return r.git("show", fmt.Sprintf("%s:%s", commit, filePath))
}

// ListFiles lists all files in a commit.
func (r *Repo) ListFiles(commit string) ([]string, error) {
	output, err := r.git("ls-tree", "-r", "--name-only", commit)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, line := range strings.Split(output, "\n") {
		if line != "" {
			files = append(files, line)
		}
	}
	return files, nil
}

// GetCommit gets a commit by ID or reference.
func (r *Repo) GetCommit(ref string) (*Commit, error) {
	commits, err := r.Log(&LogOptions{From: ref + "^", Limit: 1})
	if err != nil {
		return nil, err
	}
	if len(commits) == 0 {
		return nil, fmt.Errorf("commit not found: %s", ref)
	}
	return &commits[0], nil
}

// GetCurrentUser returns the current user from config.
func (r *Repo) GetCurrentUser() (*Signature, error) {
	name, _ := r.GetConfig("user.name")
	if name == "" {
		name = "GitForge User"
	}
	email, _ := r.GetConfig("user.email")
	if email == "" {
		email = "user@gitforge.local"
	}
	return &Signature{Name: name, Email: email}, nil
}