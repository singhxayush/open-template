package config

import (
	"flag"
	"fmt"
	"os"
)

// Config holds the application configuration.
type Config struct {
	Help    bool
	Depth   int
	Verbose bool
	Command string // Parsed command
}

// LoadConfig loads configuration from flags and potentially config files (later).
func LoadConfig() (*Config, error) {
	cfg := &Config{}

	flags := flag.NewFlagSet("template-manager", flag.ContinueOnError)

	// Move flag setup to flag.go
	SetupFlags(flags, cfg)
	if err := flags.Parse(os.Args[1:]); err != nil && err != flag.ErrHelp {
		return nil, err
	}

	// Get non-flag arguments (commands)
	args := flag.Args()
	if len(args) > 0 {
		command := args[0]
		switch command {
		case "auth", "sync", "status":
			cfg.Command = command
		default:
			fmt.Printf("Unknown command: %s\n", command)
			fmt.Println("Use '--help' to see available commands.")
			os.Exit(1)
		}
	}

	// Handle help flag
	if cfg.Help {
		PrintHelp()              // Moved help printing here and updated to use flagset
		return nil, flag.ErrHelp // Return error to signal help was shown, main can exit.
	}

	return cfg, nil
}
