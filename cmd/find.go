package cmd

import "github.com/spf13/cobra"

var findCmd = &cobra.Command{
	Use:   "find",
	Short: "Find stale resources",
	Long:  `Find stale resources`,
}

var findResourceCmd = &cobra.Command{
	Use:   "resource",
	Short: "Find stale resources",
	Long:  `Find stale resources`,
	Args:  cobra.MinimumNArgs(1),
}
