package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var RootCmd = &cobra.Command{Use: "app"}

var findStaleCmd = &cobra.Command{
	Use:   "find [resource-type] [resource-name]",
	Short: "finds all stale resources",
	Long:  "finds all stale resources",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("argument 1", args[0], args[1])
		fmt.Println("find-stale")
	},
	DisableFlagsInUseLine: true,
}

var deleteStaleCmd = &cobra.Command{
	Use:   "delete-stale",
	Short: "delete-stale",
	Long:  `delete-stale`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("delete-stale")
	},
}

func Execute() {
	RootCmd.AddCommand(findStaleCmd, deleteStaleCmd)
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
