package utils

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	treeWayConnector = "├╼ "

	twoWayConnector = "└╼ "
)

var (
	emptyMsg = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF8FA3")).Render("E\nM\nP\nT\nY")
)

// GetFileTree returns a full recursive tree view of the given directory.
// It accepts maxDepth parameter for the depth limit of the tree.
// If maxDepth is negative, the directory is traversed fully.
func GetFileTree(dir string, maxDepth int) string {
	// String builder variable to store the store the formated result or error at every-step
	var sb strings.Builder

	// Recursive function on directories
	var walk func(path string, prefix string, depth int)

	walk = func(path string, prefix string, depth int) {
		entries, err := os.ReadDir(path)
		if err != nil {
			sb.WriteString(prefix + fmt.Sprintf("Error reading directory: %v\n", err))
			return
		}

		// Slices to store directories and files separately
		var dirs []fs.DirEntry
		var files []fs.DirEntry

		// Iterate over the entires of the current path
		// to check and divide directories and files in respective slices
		for _, entry := range entries {
			if entry.IsDir() {
				dirs = append(dirs, entry)
			} else {
				files = append(files, entry)
			}
		}

		// Sort directories and files: directories first (alphabetically), then files.
		sort.Slice(dirs, func(i, j int) bool {
			return strings.ToLower(dirs[i].Name()) < strings.ToLower(dirs[j].Name())
		})
		sort.Slice(files, func(i, j int) bool {
			return strings.ToLower(files[i].Name()) < strings.ToLower(files[j].Name())
		})

		// Slice to store group-sorted dirctories and files
		allEntries := append(dirs, files...)

		for i, entry := range allEntries {
			connector := "├╼ "
			if i == len(allEntries)-1 {
				connector = "└╼ "
			}
			name := entry.Name()
			if entry.IsDir() {
				// Apply color for directories
				name = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFB2BF")).Bold(true).Render(name)
				name = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFB2BF")).Render(" ") + name
			} else {
				// Apply color for files
				name = lipgloss.NewStyle().Foreground(lipgloss.Color("#AAABB7")).Render(name)
			}
			sb.WriteString(prefix + connector + name + "\n")
			// Recurse into directories if we haven't reached maxDepth (if one is set)
			if entry.IsDir() && (maxDepth < 0 || depth < maxDepth) {
				newPrefix := prefix
				if i == len(allEntries)-1 {
					newPrefix += "   "
				} else {
					newPrefix += "│  "
				}
				walk(filepath.Join(path, entry.Name()), newPrefix, depth+1)
			}
		}
	}

	// Recursive call
	walk(dir, "", 1)
	if len(sb.String()) == 0 {
		// Handle if root directory of the template is empty
		emptyMsg := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF8FA3")).Margin(5, 4).Render("E\nM\nP\nT\nY")
		sb.WriteString(emptyMsg)
	}

	return sb.String()
}
