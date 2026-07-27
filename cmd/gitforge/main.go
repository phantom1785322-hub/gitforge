// GitForge - A beautiful, powerful Git client for the terminal and web.
package main

import (
	"fmt"
	"os"

	"github.com/gitforge/gitforge/cmd/gitforge/commands"
	"github.com/spf13/cobra"
)

var version = "0.1.0-dev"

func main() {
	rootCmd := &cobra.Command{
		Use:     "gitforge",
		Short:   "GitForge - A beautiful, powerful Git client for terminal and web",
		Long:    `GitForge is a beautiful, powerful Git client that works in your terminal and browser.
It provides visual commit graphs, interactive staging, conflict resolution,
AI-assisted commits, and seamless GitHub/GitLab integration.`,
		Version: version,
	}

	// Add subcommands
	rootCmd.AddCommand(commands.NewTUICmd())
	rootCmd.AddCommand(commands.NewWebCmd())
	rootCmd.AddCommand(commands.NewInitCmd())
	rootCmd.AddCommand(commands.NewCloneCmd())
	rootCmd.AddCommand(commands.NewStatusCmd())
	rootCmd.AddCommand(commands.NewLogCmd())
	rootCmd.AddCommand(commands.NewBranchCmd())
	rootCmd.AddCommand(commands.NewCommitCmd())
	rootCmd.AddCommand(commands.NewDiffCmd())
	rootCmd.AddCommand(commands.NewRemoteCmd())
	rootCmd.AddCommand(commands.NewStashCmd())
	rootCmd.AddCommand(commands.NewTagCmd())
	rootCmd.AddCommand(commands.NewConfigCmd())
	rootCmd.AddCommand(commands.NewPluginCmd())
	rootCmd.AddCommand(commands.NewDoctorCmd())
	rootCmd.AddCommand(commands.NewVersionCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}