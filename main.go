// main.go
package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	style "template-manager/internal/ui/style"
	"template-manager/utils"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Absolute route to your template folder.
const templateDir string = "/Users/ayushkumar/programming/templates"

// ----- Application Stages -----
const (
	stageSelectTemplate = iota
	stageProjectName
	stageCopying
	stageDone
)

// blinkMsg is sent periodically to toggle the blink state.
type blinkMsg struct{}

// ----- Data Types -----
// op represents a file or directory operation.
type op struct {
	opType  string // "mkdir" or "copy"
	relPath string
}

// model holds the application state.
type model struct {
	stage int

	// Stage 0: Template selection.
	templates []string
	cursor    int

	// Search-related fields for template selection.
	searchMode    bool
	searchQuery   string
	searchResults []string
	searchCursor  int

	// Stage 1: Project name input.
	inputBuffer string
	projectName string

	// Stage 2: Copying process.
	ops            []op // list of operations to perform
	currentOpIndex int

	// A single log message - only one log appears at a time.
	currentLog string

	// Spinner used during copying.
	spinner spinner.Model

	// Paths for copying.
	sourceDir string // full path of the selected template
	destDir   string // destination directory (created in CWD)

	// Tree depth parameter; negative means unlimited.
	treeDepth int

	// Blink state for cursor.
	blink bool

	// Show help instructions panel.
	showHelp bool

	// Error (if any).
	err error
}

// Messages for the copying process.
type opProcessedMsg struct {
	op  op
	err error
}

type copyFinishedMsg struct{}

// ----- Helper Functions -----

// clipText clips the given text to a maximum number of lines.
func clipText(text string, maxLines int) string {
	lines := strings.Split(text, "\n")
	if len(lines) > maxLines {
		lines = lines[:maxLines]
		lines = append(lines, "...")
	}
	return strings.Join(lines, "\n")
}

// copyFile copies a file from src to dst.
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	// Ensure the destination directory exists.
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

// buildOps walks the source directory and builds a list of operations.
func buildOps(source string) ([]op, error) {
	var ops []op
	err := filepath.Walk(source, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}

		if info.IsDir() {
			ops = append(ops, op{opType: "mkdir", relPath: rel})
		} else {
			ops = append(ops, op{opType: "copy", relPath: rel})
		}
		return nil
	})
	return ops, err
}

// nextOpCmd processes the next operation.
func nextOpCmd(m model) tea.Cmd {
	return func() tea.Msg {
		if m.currentOpIndex >= len(m.ops) {
			return copyFinishedMsg{}
		}
		currentOp := m.ops[m.currentOpIndex]
		var err error
		switch currentOp.opType {
		case "mkdir":
			destPath := filepath.Join(m.destDir, currentOp.relPath)
			err = os.MkdirAll(destPath, 0755)
		case "copy":
			srcPath := filepath.Join(m.sourceDir, currentOp.relPath)
			destPath := filepath.Join(m.destDir, currentOp.relPath)
			err = copyFile(srcPath, destPath)
		}
		// Simulate a slight delay.
		time.Sleep(200 * time.Millisecond)
		return opProcessedMsg{op: currentOp, err: err}
	}
}

// loadTemplates returns a slice of template names from the templates directory.
func loadTemplates(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var templates []string
	for _, entry := range entries {
		if entry.IsDir() {
			templates = append(templates, entry.Name())
		}
	}
	return templates, nil
}

// Command suggestion style (dimmed).
var commandStyle = lipgloss.NewStyle().Faint(true)

// ----- Bubble Tea Model Methods -----
func initialModel() model {
	templates, err := loadTemplates(templateDir)
	if err != nil || len(templates) == 0 {
		fmt.Println("Error loading templates or no templates found in", templateDir)
		os.Exit(1)
	}

	// Initialize the spinner with the Jump spinner.
	s := spinner.New()
	s.Spinner = spinner.Jump
	// Adjust the FPS if desired.
	spinner.Jump.FPS = 12
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFA500")).Margin(0, 0)

	return model{
		stage:      stageSelectTemplate,
		templates:  templates,
		spinner:    s,
		treeDepth:  -1, // unlimited depth by default; can be updated via flag.
		searchMode: false,
		blink:      true,
		showHelp:   false,
	}
}

func (m model) Init() tea.Cmd {
	// Start the blink ticker.
	return tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg { return blinkMsg{} })
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// Global: Exit immediately if Ctrl+C is pressed.
	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "ctrl+c" {
		return m, tea.Quit
	}

	// Check for the help toggle key "?"
	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "?" {
		m.showHelp = !m.showHelp
		return m, nil
	}

	switch msg := msg.(type) {
	case blinkMsg:
		m.blink = !m.blink
		cmds = append(cmds, tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg {
			return blinkMsg{}
		}))

	// ----- Stage 0: Template Selection -----
	case tea.KeyMsg:
		// Process keys for search mode if active.
		if m.stage == stageSelectTemplate && m.searchMode {
			switch msg.Type {
			case tea.KeyEsc:
				// Exit search mode and clear search state.
				m.searchMode = false
				m.searchQuery = ""
				m.searchResults = nil
				m.searchCursor = 0
			case tea.KeyEnter:
				// If there are any suggestions, select the current suggestion and move to the next stage.
				if len(m.searchResults) > 0 {
					selection := m.searchResults[m.searchCursor]
					// Find the index of the selection in the full list.
					for i, tmpl := range m.templates {
						if tmpl == selection {
							m.cursor = i
							break
						}
					}
					m.sourceDir = filepath.Join(templateDir, selection)
					m.stage = stageProjectName
				}
				m.searchMode = false
				m.searchQuery = ""
				m.searchResults = nil
				m.searchCursor = 0
			case tea.KeyBackspace:
				if len(m.searchQuery) > 0 {
					m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
				}
				// Recalculate search results.
				m.searchResults = utils.FilterTemplates(m.templates, m.searchQuery)
				if m.searchCursor >= len(m.searchResults) && len(m.searchResults) > 0 {
					m.searchCursor = len(m.searchResults) - 1
				}
			case tea.KeyUp, tea.KeyShiftUp:
				if m.searchCursor > 0 {
					m.searchCursor--
				}
			case tea.KeyDown, tea.KeyShiftDown:
				if m.searchCursor < len(m.searchResults)-1 {
					m.searchCursor++
				}
			default:
				// Append any other character to the search query.
				m.searchQuery += msg.String()
				m.searchResults = utils.FilterTemplates(m.templates, m.searchQuery)
				m.searchCursor = 0
			}
			return m, tea.Batch(cmds...)
		}

		// Normal key handling (outside of search mode)
		if m.stage == stageSelectTemplate {
			switch msg.String() {
			case "/":
				// Enter search mode.
				m.searchMode = true
				m.searchQuery = ""
				m.searchResults = m.templates // show all initially
				m.searchCursor = 0
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.templates)-1 {
					m.cursor++
				}
			case "enter":
				// When a template is selected, set the source directory.
				selectedTemplate := m.templates[m.cursor]
				m.sourceDir = filepath.Join(templateDir, selectedTemplate)
				// Transition to project name input.
				m.stage = stageProjectName
			case "q":
				return m, tea.Quit
			}
		} else if m.stage == stageProjectName {
			switch msg.Type {
			case tea.KeyEnter:
				m.projectName = strings.TrimSpace(m.inputBuffer)
				if m.projectName == "" {
					// Do nothing if project name is empty.
					return m, nil
				}
				// Create the destination directory in the current working directory.
				cwd, err := os.Getwd()
				if err != nil {
					m.err = fmt.Errorf("Error getting CWD: %v", err)
					return m, tea.Quit
				}
				m.destDir = filepath.Join(cwd, m.projectName)
				if err := os.Mkdir(m.destDir, 0755); err != nil {
					m.err = fmt.Errorf("Error creating project directory: %v", err)
					return m, tea.Quit
				}
				// Build copy operations.
				ops, err := buildOps(m.sourceDir)
				if err != nil {
					m.err = fmt.Errorf("Error building copy operations: %v", err)
					return m, tea.Quit
				}
				m.ops = ops
				m.currentOpIndex = 0
				m.currentLog = ""
				m.stage = stageCopying
				// Begin processing copy operations and start spinner ticking.
				cmds = append(cmds, nextOpCmd(m), m.spinner.Tick)
			case tea.KeyBackspace:
				if len(m.inputBuffer) > 0 {
					m.inputBuffer = m.inputBuffer[:len(m.inputBuffer)-1]
				}
			default:
				// Append typed characters.
				m.inputBuffer += msg.String()
			}
		}

	// ----- Stage 2: Copying Process -----
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	case opProcessedMsg:
		// Update the single log line.
		if msg.op.opType == "mkdir" {
			m.currentLog = fmt.Sprintf("Created directory: %s", msg.op.relPath)
		} else {
			m.currentLog = fmt.Sprintf("Copied file: %s", msg.op.relPath)
		}
		if msg.err != nil {
			m.currentLog += fmt.Sprintf(" [Error: %v]", msg.err)
		}
		m.currentOpIndex++
		if m.currentOpIndex < len(m.ops) {
			cmds = append(cmds, nextOpCmd(m))
		} else {
			cmds = append(cmds, func() tea.Msg { return copyFinishedMsg{} })
		}
	case copyFinishedMsg:
		m.currentLog = "Project " + "\"" + m.projectName + "\"" + " created successfully!"
		m.stage = stageDone
	}

	// ----- Stage 3: Done -----
	if m.stage == stageDone {
		return m, tea.Quit
	}

	return m, tea.Batch(cmds...)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (m model) View() string {
	if m.err != nil {
		return style.DocStyle.Render(style.ErrorStyle.Render(fmt.Sprintf("Error: %v", m.err)))
	}

	var body string

	switch m.stage {
	case stageSelectTemplate:
		var leftPanel string
		// Build left panel contents.
		if m.searchMode {
			var sb strings.Builder
			// Append search prompt with blinking cursor.
			cursor := ""
			if m.blink {
				cursor = style.CursorStyle.Render("|")
			}
			sb.WriteString("Search: " + m.searchQuery + cursor + "\n\n")
			// If there are no results, show a message.
			if len(m.searchResults) == 0 {
				sb.WriteString("No matching templates")
			} else {
				// Show suggestions with dynamic highlighting.
				for i, tmpl := range m.searchResults {
					curs := lipgloss.NewStyle().Foreground(lipgloss.Color("#D2F8B0")).Render("⬥")
					itemStyle := style.ListItemStyle
					if m.searchCursor == i {
						curs = "⬥"
						itemStyle = lipgloss.NewStyle().
							Foreground(lipgloss.Color("#A6E3A1")).
							MarginLeft(2).Bold(true)
					}
					sb.WriteString(fmt.Sprintf("%s%s\n", curs, itemStyle.Render(tmpl)))
				}
			}
			leftPanel = style.LeftPanelStyle.Render(sb.String())
		} else {
			// Normal list view.
			visibleLines := 7 // Adjust based on space after prompt/instructions.
			start := m.cursor - visibleLines/2
			if start < 0 {
				start = 0
			}
			end := start + visibleLines
			if end > len(m.templates) {
				end = len(m.templates)
				start = max(0, end-visibleLines)
			}
			var listBuilder strings.Builder
			for i := start; i < end; i++ {
				tmpl := m.templates[i]
				curs := "⬦"
				itemStyle := style.ListItemStyle
				if m.cursor == i {
					curs = lipgloss.NewStyle().Foreground(lipgloss.Color("#D2F8B0")).Render("⬥")
					itemStyle = lipgloss.NewStyle().
						Foreground(lipgloss.Color("#A6E3A1")).
						MarginLeft(2).Bold(true)
				}
				listBuilder.WriteString(fmt.Sprintf("%s%s\n", curs, itemStyle.Render(tmpl)))
			}
			leftPanel = style.LeftPanelStyle.Render(listBuilder.String())
		}

		// Append fixed command instructions below the left panel.
		var instructions string
		if m.searchMode {
			instructions = " Navigate : ↑/↓ or j/k\n Select   : Enter\n ESC      : Cancel Search\n Exit     : ctrl+c"
		} else if m.showHelp {
			instructions = "Navigate: ↑/↓ or j/k\n" +
				"Select  : Enter\n" +
				"Search  : /\n" +
				"Help    : ?\n" +
				"Exit    : ctrl+c"
		} else {
			instructions = lipgloss.NewStyle().
				// Border(lipgloss.NormalBorder()).
				PaddingLeft(5).PaddingRight(6).
				PaddingTop(1).PaddingBottom(1).
				Render("/ find • q quit • ? help")
		}

		instructions = commandStyle.Render(instructions)
		leftContent := lipgloss.JoinVertical(lipgloss.Left, leftPanel, instructions)

		// For the right panel, show the tree of the currently highlighted template.
		var selectedTemplate string
		if m.searchMode && len(m.searchResults) > 0 {
			selectedTemplate = m.searchResults[m.searchCursor]
		} else {
			selectedTemplate = m.templates[m.cursor]
		}
		templatePath := filepath.Join(templateDir, selectedTemplate)
		rightContent := utils.GetFileTree(templatePath, m.treeDepth)
		rightPanel := style.RightPanelStyle.Render(rightContent)

		// Horizontally join the fixed left content with the right panel.
		body = lipgloss.JoinHorizontal(lipgloss.Top, leftContent, rightPanel)

	case stageProjectName:
		// Show the project name prompt with blinking cursor.
		cursor := ""
		if m.blink {
			cursor = style.CursorStyle.Render("|")
		}
		body = fmt.Sprintf("Enter project name: %s%s\n\nPress Ctrl+C to exit at any point.", m.inputBuffer, cursor)

	case stageCopying:
		// Render a single log line with the spinner.
		body = fmt.Sprintf("%s %s\n\nPress Ctrl+C to exit at any point.", m.spinner.View(), m.currentLog)

	case stageDone:
		body = fmt.Sprint("Done")
	}

	return style.DocStyle.Render(fmt.Sprintf("%s\n%s\n", style.HeaderStyle.Render("Template Manager ⚡"), body))
}

// ----- Main -----
func main() {
	// Parse flags and commands correctly
	cf := utils.ParseFlags()
	cp := utils.ParseCommands()

	// Execute commands (if any)
	utils.Execute(cf, cp)

	// If a command was executed, exit before launching UI
	if cp.Command != "" {
		os.Exit(0)
	}

	// Initialize UI model
	m := initialModel()
	m.treeDepth = cf.Depth

	// Run Bubble Tea program
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
