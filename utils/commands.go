package utils

import (
	"flag"
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
)

var (
	headlineStyles = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FCF8F7")).
			Margin(1)

	commandStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A6E3A1")).
			MarginLeft(4)

	descriptionStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("##FFFFF0"))
)

// CmdFlags - Struct to hold CLI flags
// for example --help or -h help
type CmdFlags struct {
	Help    bool
	Depth   int
	Verbose bool
}

// CmdParams - Struct to hold CLI commands
type CmdParams struct {
	Command string
}

// ParseFlags - Parses CLI flags
func ParseFlags() *CmdFlags {
	cf := &CmdFlags{}

	// Define flags
	flag.BoolVar(&cf.Help, "help", false, "Show help message")
	flag.IntVar(&cf.Depth, "depth", 1, "Set max depth for file tree (-1 for unlimited)")
	flag.BoolVar(&cf.Verbose, "verbose", false, "Enable verbose logging")

	// Parse known flags
	flag.Parse()

	return cf
}

// ParseCommands - Detects a valid command (ignoring flags)
func ParseCommands() *CmdParams {
	cp := &CmdParams{}

	// Get non-flag arguments (commands)
	args := flag.Args()
	if len(args) > 0 {
		command := args[0]
		switch command {
		case "auth", "sync", "status":
			cp.Command = command
		default:
			fmt.Printf("Unknown command: %s\n", command)
			fmt.Println("Use '--help' to see available commands.")
			os.Exit(1)
		}
	}

	return cp
}

// Execute - Runs commands and flags
func Execute(cf *CmdFlags, cp *CmdParams) {
	// Handle flags first
	if cf.Help {
		PrintHelp()
		os.Exit(0)
	}

	// Handle commands
	switch cp.Command {
	case "auth":
		fmt.Println("Initializing authentication...")
	case "sync":
		fmt.Println("Syncing cloud changes with local machine...")
	case "status":
		fmt.Println("Checking system status...")
	}
}

// PrintHelp - Displays help menu
func PrintHelp() {
	fmt.Println(headlineStyles.MarginTop(0).Render("USAGE:"))
	fmt.Printf("  [command] [flags]\n")

	fmt.Printf("%v\n", headlineStyles.Render("Commands:"))
	fmt.Printf("%v\t%v\n", commandStyle.Render("auth"), descriptionStyle.Render("Initialize authentication"))
	fmt.Printf("%v\t%v\n", commandStyle.Render("sync"), descriptionStyle.Render("Sync cloud changes on the local machine"))
	fmt.Printf("%v\t%v\n", commandStyle.Render("status"), descriptionStyle.Render("Show system and sync status"))

	fmt.Println(headlineStyles.Render("Flags:"))
	fmt.Printf("%v\t%v\n", commandStyle.Render("--help"), descriptionStyle.Render("Show this help message"))
	fmt.Printf("%v\t%v\n", commandStyle.Render("--depth n"), descriptionStyle.Render("Set depth of file tree visualization (-1 for unlimited)"))
	fmt.Printf("%v\t%v\n", commandStyle.Render("--config"), descriptionStyle.Render("Specify config file path"))
	fmt.Printf("%v\t%v\n", commandStyle.Render("--verbose"), descriptionStyle.Render("Enable verbose logging"))

	fmt.Println(headlineStyles.Render("Examples:"))
	fmt.Println("  go run main.go auth")
	fmt.Println("  go run main.go sync")
	fmt.Println("  go run main.go status")
	fmt.Println("  go run main.go --depth=2 --verbose")

}
