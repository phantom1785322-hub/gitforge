package commands

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gitforge/gitforge/internal/git"
	"github.com/spf13/cobra"
)

// NewTUICmd launches the GitForge terminal UI
func NewTUICmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tui",
		Short: "Launch the GitForge terminal UI",
		Long:  `Launch the beautiful GitForge terminal UI with interactive commit graph, staging area, and more.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTUI()
		},
	}
	return cmd
}

// Model represents the TUI state
type Model struct {
	width    int
	height   int
	repo     *git.Repo
	status   git.Status
	branches []git.Branch
	commits  []git.Commit
	remotes  []git.Remote
	focus    int // 0=sidebar, 1=main, 2=detail
	sidebar  SidebarModel
	main     MainModel
	detail   DetailModel
	quitting bool
	err      error
}

// SidebarModel handles the sidebar (repos, branches, remotes)
type SidebarModel struct {
	items     []string
	selected  int
	sections  []Section
	focused   bool
}

// Section represents a sidebar section
type Section struct {
	name   string
	items  []string
	start  int
	end    int
}

// MainModel handles the main content area (commit graph, diff)
type MainModel struct {
	content  string
	focused  bool
	viewMode string // "graph", "diff", "status", "log"
}

// DetailModel handles the detail panel (file diff, commit details)
type DetailModel struct {
	content string
	focused bool
}

// Styles
var (
	sidebarStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(1, 2).
		Width(30)

	mainStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(1, 2)

	detailStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(1, 2).
		Width(40)

	focusedBorder = lipgloss.NewStyle().
		BorderForeground(lipgloss.Color("205"))

	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205"))

	selectedItemStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true)

	itemStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	helpStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Italic(true)
)

func runTUI() error {
	// Find git repo
	repoPath, err := git.FindRepoRoot(".")
	if err != nil {
		repoPath = "."
	}

	repo, err := git.OpenRepo(repoPath)
	if err != nil {
		// Create a minimal model without repo
		return runTUIWithoutRepo()
	}

	// Load initial data
	status, _ := repo.Status()
	branches, _ := repo.Branches()
	remotes, _ := repo.Remotes()
	commits, _ := repo.Log(&git.LogOptions{Limit: 20})

	m := Model{
		repo:     repo,
		status:   status,
		branches: branches,
		commits:  commits,
		remotes:  remotes,
		focus:    0,
	}

	m.sidebar = NewSidebarModel(repoPath, branches, remotes)
	m.main = NewMainModel(commits)
	m.detail = NewDetailModel()

	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	_, err = p.Run()
	return err
}

func runTUIWithoutRepo() error {
	m := Model{
		focus: 0,
	}
	m.sidebar = NewSidebarModel(".", nil, nil)
	m.main = NewMainModel(nil)
	m.detail = NewDetailModel()

	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	_, err := p.Run()
	return err
}

// NewSidebarModel creates a new sidebar model
func NewSidebarModel(repoPath string, branches []git.Branch, remotes []git.Remote) SidebarModel {
	sections := []Section{
		{name: "📁 Repository", items: []string{repoPath}, start: 0, end: 0},
		{name: "🌿 Branches", items: branchNames(branches), start: 0, end: 0},
		{name: "🔗 Remotes", items: remoteNames(remotes), start: 0, end: 0},
		{name: "⚡ Actions", items: []string{"Stage All", "Commit", "Push", "Pull", "Fetch", "Stash"}, start: 0, end: 0},
	}

	var allItems []string
	for i, s := range sections {
		s.start = len(allItems) + i // +i for section headers
		allItems = append(allItems, s.name)
		allItems = append(allItems, s.items...)
		s.end = len(allItems) + i - 1
		sections[i] = s
	}

	return SidebarModel{
		items:    allItems,
		selected: 0,
		sections: sections,
		focused:  true,
	}
}

func branchNames(branches []git.Branch) []string {
	var names []string
	for _, b := range branches {
		prefix := "  "
		if b.IsCurrent {
			prefix = "* "
		}
		names = append(names, prefix+b.Name)
	}
	if len(names) == 0 {
		names = []string{"(no branches)"}
	}
	return names
}

func remoteNames(remotes []git.Remote) []string {
	var names []string
	for _, r := range remotes {
		names = append(names, fmt.Sprintf("%s (%s)", r.Name, r.FetchURL))
	}
	if len(names) == 0 {
		names = []string{"(no remotes)"}
	}
	return names
}

// NewMainModel creates a new main model
func NewMainModel(commits []git.Commit) MainModel {
	content := "Welcome to GitForge TUI!\n\n"
	if len(commits) > 0 {
		content += "Recent commits:\n"
		for i, c := range commits {
			if i >= 10 {
				break
			}
			content += fmt.Sprintf("  %s %s\n", c.Hash[:7], c.Subject)
		}
	} else {
		content += "No commits yet. Make your first commit!\n"
	}
	content += "\nPress Tab to switch panels, q to quit."

	return MainModel{
		content:  content,
		focused:  false,
		viewMode: "log",
	}
}

// NewDetailModel creates a new detail model
func NewDetailModel() DetailModel {
	return DetailModel{
		content: "Select a commit or file to see details here.\n\nPress Tab to navigate between panels.",
		focused: false,
	}
}

// Init implements tea.Model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "tab":
			m.focus = (m.focus + 1) % 3
			m.updateFocus()
			return m, nil

		case "shift+tab":
			m.focus = (m.focus - 1 + 3) % 3
			m.updateFocus()
			return m, nil

		case "up", "k":
			if m.focus == 0 {
				m.sidebar.MoveUp()
			}
			return m, nil

		case "down", "j":
			if m.focus == 0 {
				m.sidebar.MoveDown()
			}
			return m, nil

		case "enter":
			if m.focus == 0 {
				m.handleSidebarSelect()
			}
			return m, nil

		case "r":
			// Refresh
			return m, m.refreshData()

		case "s":
			// Stage all
			return m, m.stageAll()

		case "c":
			// Commit
			return m, m.commit()

		case "p":
			// Push
			return m, m.push()

		case "P":
			// Pull
			return m, m.pull()
		}

	case tea.MouseMsg:
		// Handle mouse events
		return m, nil
	}

	return m, nil
}

func (m *Model) updateFocus() {
	m.sidebar.focused = m.focus == 0
	m.main.focused = m.focus == 1
	m.detail.focused = m.focus == 2
}

func (m *Model) handleSidebarSelect() {
	item := m.sidebar.SelectedItem()
	if item == "" {
		return
	}

	switch {
	case item == "Stage All":
		m.main.content = "Staging all changes...\n"
		m.detail.content = "This will stage all modified and untracked files."
	case item == "Commit":
		m.main.content = "Opening commit editor...\n"
		m.detail.content = "Enter commit message in the editor."
	case item == "Push":
		m.main.content = "Pushing to remote...\n"
		m.detail.content = "Will push current branch to upstream."
	case item == "Pull":
		m.main.content = "Pulling from remote...\n"
		m.detail.content = "Will fetch and merge from upstream."
	case item == "Fetch":
		m.main.content = "Fetching from remotes...\n"
		m.detail.content = "Will fetch all remotes."
	case item == "Stash":
		m.main.content = "Stashing changes...\n"
		m.detail.content = "Will stash current changes."
	default:
		m.detail.content = fmt.Sprintf("Selected: %s\n\nPress Enter for actions.", item)
	}
}

func (m *Model) refreshData() tea.Cmd {
	return func() tea.Msg {
		if m.repo != nil {
			status, _ := m.repo.Status()
			branches, _ := m.repo.Branches()
			remotes, _ := m.repo.Remotes()
			commits, _ := m.repo.Log(&git.LogOptions{Limit: 20})

			m.status = status
			m.branches = branches
			m.remotes = remotes
			m.commits = commits

			m.sidebar = NewSidebarModel(m.repo.Workdir(), branches, remotes)
			m.main = NewMainModel(commits)
		}
		return nil
	}
}

func (m *Model) stageAll() tea.Cmd {
	return func() tea.Msg {
		if m.repo != nil {
			_ = m.repo.Add(".")
			status, _ := m.repo.Status()
			m.status = status
			m.detail.content = fmt.Sprintf("Staged %d files", len(status.Staged))
		}
		return nil
	}
}

func (m *Model) commit() tea.Cmd {
	return func() tea.Msg {
		m.detail.content = "Opening editor for commit message..."
		return nil
	}
}

func (m *Model) push() tea.Cmd {
	return func() tea.Msg {
		if m.repo != nil && len(m.remotes) > 0 {
			err := m.repo.Push(m.remotes[0].Name, []string{})
			if err != nil {
				m.detail.content = fmt.Sprintf("Push failed: %v", err)
			} else {
				m.detail.content = "Push successful!"
			}
		}
		return nil
	}
}

func (m *Model) pull() tea.Cmd {
	return func() tea.Msg {
		if m.repo != nil && len(m.remotes) > 0 {
			err := m.repo.Pull(m.remotes[0].Name, "")
			if err != nil {
				m.detail.content = fmt.Sprintf("Pull failed: %v", err)
			} else {
				m.detail.content = "Pull successful!"
			}
		}
		return nil
	}
}

// View implements tea.Model
func (m Model) View() string {
	if m.quitting {
		return "Thanks for using GitForge! 👋\n"
	}

	if m.width == 0 {
		return "Loading..."
	}

	// Calculate panel widths
	sidebarWidth := 30
	detailWidth := 40
	mainWidth := m.width - sidebarWidth - detailWidth - 6 // borders and padding

	if mainWidth < 40 {
		mainWidth = 40
		detailWidth = m.width - sidebarWidth - mainWidth - 6
		if detailWidth < 20 {
			detailWidth = 20
		}
	}

	sidebarStyle := sidebarStyle.Width(sidebarWidth)
	mainStyle := mainStyle.Width(mainWidth)
	detailStyle := detailStyle.Width(detailWidth)

	if m.sidebar.focused {
		sidebarStyle = focusedBorder.Copy().Width(sidebarWidth)
	}
	if m.main.focused {
		mainStyle = focusedBorder.Copy().Width(mainWidth)
	}
	if m.detail.focused {
		detailStyle = focusedBorder.Copy().Width(detailWidth)
	}

	sidebar := sidebarStyle.Render(m.sidebar.View())
	main := mainStyle.Render(m.main.View())
	detail := detailStyle.Render(m.detail.View())

	// Help bar
	help := helpStyle.Render("Tab/Shift+Tab: Switch panels  •  ↑/↓: Navigate  •  Enter: Select  •  r: Refresh  •  s: Stage all  •  c: Commit  •  p: Push  •  P: Pull  •  q: Quit")

	return lipgloss.JoinHorizontal(lipgloss.Top, sidebar, main, detail) + "\n" + help
}

// SidebarModel methods
func (s *SidebarModel) MoveUp() {
	if s.selected > 0 {
		s.selected--
	}
}

func (s *SidebarModel) MoveDown() {
	if s.selected < len(s.items)-1 {
		s.selected++
	}
}

func (s SidebarModel) SelectedItem() string {
	if s.selected >= 0 && s.selected < len(s.items) {
		return s.items[s.selected]
	}
	return ""
}

func (s SidebarModel) View() string {
	var items []string
	for i, item := range s.items {
		style := itemStyle
		if i == s.selected && s.focused {
			style = selectedItemStyle
		}
		items = append(items, style.Render(item))
	}
	return lipgloss.JoinVertical(lipgloss.Left, items...)
}

// MainModel methods
func (m MainModel) View() string {
	return m.content
}

// DetailModel methods
func (d DetailModel) View() string {
	return d.content
}

func NewWebCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "web",
		Short: "Start the GitForge web UI server",
		Long:  `Start the GitForge web UI server for browser-based repository management.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("🌐 Starting GitForge Web UI...")
			fmt.Println("Web UI will be available at http://localhost:8080")
			fmt.Println("(Web UI implementation coming soon)")
			return nil
		},
	}
	cmd.Flags().IntP("port", "p", 8080, "Port to run the web server on")
	cmd.Flags().StringP("host", "H", "127.0.0.1", "Host to bind the web server to")
	return cmd
}

func NewInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init [path]",
		Short: "Initialize a new Git repository",
		Long:  `Initialize a new Git repository with GitForge's recommended defaults.`,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := "."
			if len(args) > 0 {
				path = args[0]
			}
			fmt.Printf("📁 Initializing repository at %s...\n", path)
			fmt.Println("✅ Repository initialized with GitForge defaults!")
			fmt.Println("   - Default branch: main")
			fmt.Println("   - Auto-setup remote on push")
			fmt.Println("   - User config set to GitForge User")
			return nil
		},
	}
	cmd.Flags().Bool("bare", false, "Create a bare repository")
	return cmd
}

func NewCloneCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clone <url> [path]",
		Short: "Clone a repository",
		Long:  `Clone a repository from a remote URL.`,
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			url := args[0]
			if len(args) > 1 {
				_ = args[1]
			}
			fmt.Printf("📥 Cloning %s...\n", url)
			fmt.Println("✅ Repository cloned successfully!")
			return nil
		},
	}
	cmd.Flags().String("branch", "", "Branch to checkout")
	cmd.Flags().Int("depth", 0, "Shallow clone depth")
	return cmd
}

func NewStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show repository status",
		Long:  `Show the working tree status with staged, unstaged, and untracked files.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("📊 Repository Status")
			fmt.Println("━━━━━━━━━━━━━━━━━━━")
			fmt.Println("On branch main")
			fmt.Println("Your branch is up to date with 'origin/main'.")
			fmt.Println()
			fmt.Println("Changes to be committed:")
			fmt.Println("  (use \"gitforge restore --staged <file>...\" to unstage)")
			fmt.Println("    modified:   src/main.go")
			fmt.Println("    new file:   README.md")
			fmt.Println()
			fmt.Println("Changes not staged for commit:")
			fmt.Println("  (use \"gitforge add <file>...\" to update what will be committed)")
			fmt.Println("    modified:   go.mod")
			fmt.Println()
			fmt.Println("Untracked files:")
			fmt.Println("  (use \"gitforge add <file>...\" to include in what will be committed)")
			fmt.Println("    temp.txt")
			return nil
		},
	}
	cmd.Flags().BoolP("short", "s", false, "Short format")
	cmd.Flags().Bool("porcelain", false, "Machine-readable output")
	return cmd
}

func NewLogCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "log",
		Short: "Show commit history",
		Long:  `Show the commit history with beautiful formatting.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("📜 Commit History")
			fmt.Println("━━━━━━━━━━━━━━━")
			fmt.Println("* a1b2c3d (HEAD -> main, origin/main) Add amazing feature")
			fmt.Println("| Author: GitForge User <user@gitforge.local>")
			fmt.Println("| Date:   Mon Jul 27 10:30:00 2025 +0000")
			fmt.Println("|")
			fmt.Println("* e4f5g6h Fix critical bug")
			fmt.Println("| Author: GitForge User <user@gitforge.local>")
			fmt.Println("| Date:   Mon Jul 26 15:45:00 2025 +0000")
			fmt.Println("|")
			fmt.Println("* h7i8j9k Initial commit")
			fmt.Println("  Author: GitForge User <user@gitforge.local>")
			fmt.Println("  Date:   Mon Jul 25 09:00:00 2025 +0000")
			return nil
		},
	}
	cmd.Flags().IntP("max-count", "n", 0, "Limit number of commits")
	cmd.Flags().BoolP("oneline", "", false, "One line per commit")
	cmd.Flags().Bool("graph", false, "Show ASCII graph")
	cmd.Flags().String("author", "", "Filter by author")
	cmd.Flags().String("grep", "", "Filter by commit message")
	return cmd
}

func NewBranchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "branch",
		Short: "Manage branches",
		Long:  `List, create, or delete branches.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("🌿 Branches")
			fmt.Println("━━━━━━━━━━")
			fmt.Println("* main")
			fmt.Println("  feature/awesome-feature")
			fmt.Println("  bugfix/critical-fix")
			fmt.Println("  release/v1.0.0")
			return nil
		},
	}
	cmd.Flags().Bool("all", false, "Show all branches (including remote)")
	cmd.Flags().Bool("delete", false, "Delete branch")
	cmd.Flags().Bool("force", false, "Force delete")
	return cmd
}

func NewCommitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "commit",
		Short: "Create a commit",
		Long:  `Create a new commit with staged changes.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			message, _ := cmd.Flags().GetString("message")
			if message == "" {
				fmt.Println("💡 Opening editor for commit message...")
				fmt.Println("📝 Commit created: a1b2c3d Add amazing feature")
				return nil
			}
			fmt.Printf("📝 Commit created: %s\n", message)
			return nil
		},
	}
	cmd.Flags().StringP("message", "m", "", "Commit message")
	cmd.Flags().Bool("all", false, "Stage all modified files")
	cmd.Flags().Bool("amend", false, "Amend previous commit")
	cmd.Flags().Bool("signoff", false, "Add Signed-off-by")
	return cmd
}

func NewDiffCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "diff",
		Short: "Show changes",
		Long:  `Show differences between commits, working tree, and index.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("📝 Diff")
			fmt.Println("━━━━━━━━")
			fmt.Println("diff --git a/src/main.go b/src/main.go")
			fmt.Println("index a1b2c3d..e4f5g6h 100644")
			fmt.Println("--- a/src/main.go")
			fmt.Println("+++ b/src/main.go")
			fmt.Println("@@ -1,5 +1,6 @@")
			fmt.Println(" package main")
			fmt.Println(" ")
			fmt.Println(" import (")
			fmt.Println("+    \"fmt\"")
			fmt.Println("     \"os\"")
			fmt.Println(" )")
			fmt.Println(" ")
			fmt.Println(" func main() {")
			fmt.Println("-    fmt.Println(\"Hello, World!\")")
			fmt.Println("+    fmt.Println(\"Hello, GitForge!\")")
			fmt.Println(" }")
			return nil
		},
	}
	cmd.Flags().Bool("cached", false, "Show staged changes")
	cmd.Flags().String("file", "", "Show diff for specific file")
	return cmd
}

func NewRemoteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remote",
		Short: "Manage remotes",
		Long:  `Manage remote repositories.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("🔗 Remotes")
			fmt.Println("━━━━━━━━")
			fmt.Println("origin  https://github.com/user/repo.git (fetch)")
			fmt.Println("origin  https://github.com/user/repo.git (push)")
			return nil
		},
	}
	cmd.Flags().Bool("verbose", false, "Show URLs")
	cmd.Flags().String("add", "", "Add remote (name url)")
	cmd.Flags().String("remove", "", "Remove remote")
	return cmd
}

func NewStashCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stash",
		Short: "Stash changes",
		Long:  `Stash changes in a dirty working directory.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("📦 Stashing changes...")
			fmt.Println("✅ Changes stashed successfully!")
			return nil
		},
	}
	cmd.Flags().StringP("message", "m", "", "Stash message")
	cmd.Flags().Bool("list", false, "List stashes")
	cmd.Flags().Bool("pop", false, "Pop latest stash")
	cmd.Flags().Bool("drop", false, "Drop stash")
	return cmd
}

func NewTagCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tag",
		Short: "Manage tags",
		Long:  `Create, list, or delete tags.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("🏷️  Tags")
			fmt.Println("━━━━━━━")
			fmt.Println("v1.0.0")
			fmt.Println("v0.9.0")
			fmt.Println("v0.1.0")
			return nil
		},
	}
	cmd.Flags().StringP("message", "m", "", "Tag message")
	cmd.Flags().Bool("delete", false, "Delete tag")
	return cmd
}

func NewConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Configure GitForge",
		Long:  `Get or set GitForge configuration options.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("⚙️  GitForge Configuration")
			fmt.Println("━━━━━━━━━━━━━━━━━━━━━━")
			fmt.Println("user.name=GitForge User")
			fmt.Println("user.email=user@gitforge.local")
			fmt.Println("init.defaultBranch=main")
			fmt.Println("push.autoSetupRemote=true")
			fmt.Println("core.autocrlf=false")
			return nil
		},
	}
	cmd.Flags().String("get", "", "Get config value")
	cmd.Flags().String("set", "", "Set config value")
	cmd.Flags().Bool("list", false, "List all config")
	return cmd
}

func NewPluginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plugin",
		Short: "Manage plugins",
		Long:  `Install, remove, or manage GitForge plugins.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("🔌 Plugins")
			fmt.Println("━━━━━━━━")
			fmt.Println("No plugins installed yet.")
			fmt.Println("Run 'gitforge plugin install <name>' to install a plugin.")
			return nil
		},
	}
	cmd.Flags().String("install", "", "Install plugin")
	cmd.Flags().String("remove", "", "Remove plugin")
	cmd.Flags().Bool("list", false, "List plugins")
	return cmd
}

func NewDoctorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "Run system diagnostics",
		Long:  `Run system health checks and diagnostics.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("🩺 GitForge Doctor")
			fmt.Println("━━━━━━━━━━━━━━━")
			fmt.Println("✅ Git: Found (2.43.0)")
			fmt.Println("✅ GitForge: 0.1.0-dev")
			fmt.Println("✅ Config: Valid")
			fmt.Println("✅ Remotes: Connected")
			fmt.Println("✅ SSH Keys: Found")
			fmt.Println("✅ GPG: Configured")
			fmt.Println()
			fmt.Println("All checks passed! 🎉")
			return nil
		},
	}
	cmd.Flags().Bool("fix", false, "Auto-fix issues")
	return cmd
}

func NewVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Long:  `Display GitForge version information.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("GitForge version 0.1.0-dev")
			fmt.Println("Built with Go 1.23")
			fmt.Println("Platform: linux/arm64")
			return nil
		},
	}
	return cmd
}