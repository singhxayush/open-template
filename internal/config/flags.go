package config

import (
	"flag"
	"fmt"

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

// SetupFlags parses and stores it in the cgf struct defined under config.go
// config.go: LoadConfig() -> SetupFlags()
func SetupFlags(flags *flag.FlagSet, cfg *Config) {
	flags.BoolVar(&cfg.Help, "help", true, "Show help message")
	flags.IntVar(&cfg.Depth, "depth", 1, "Set max depth for file tree (-1 for unlimited)")
	flags.BoolVar(&cfg.Verbose, "verbose", false, "Enable verbose logging")
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

// ExecuteCommand executes commands based on config.
// internal/app/app.go: New() -> ExecuteCommand()
// 
func ExecuteCommand(cfg *Config) {
	switch cfg.Command {

	//Auth Case
	case "auth":
		if cfg.Verbose {
			fmt.Println("Initializing authentication...")
		}
		// ... your auth logic ...
		fmt.Printf("Implement Auth Logic MF\n")

	// Sync Case
	case "sync":
		if cfg.Verbose {
			fmt.Println("Syncing cloud changes with local machine...")
		}
		// ... your sync logic ...
		fmt.Printf("Implement SYNC Logic MF\n")

	//Status Case
	case "status":
		if cfg.Verbose {
			fmt.Println("Checking system status...")
		}
		// ... your status logic ...
		fmt.Printf("Implement Status Logic MF\n")

	// Default Case
	default:
		fmt.Println("This Will never get executed")
		// if cfg.Command != "" {
		// 	fmt.Printf("Unknown command: %s\n", cfg.Command)
		// 	fmt.Println("Use '--help' to see available commands.")
		// 	os.Exit(1)
		// }
	}
}
