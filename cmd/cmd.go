package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"kube-scourgify/pkg/controller"
	"kube-scourgify/utils"
)

var findCommand = &cobra.Command{
	Use:   "find",
	Short: "Find stale resource",
	Run: func(cmd *cobra.Command, args []string) {
		resourceName, _ := cmd.Flags().GetString(utils.RESOURCE_NAME_KEY)
		resourceVersion, _ := cmd.Flags().GetString(utils.RESOURCE_VERSION_KEY)
		resourceGroup, _ := cmd.Flags().GetString(utils.RESOURCE_GROUP_KEY)
		resourceKind, _ := cmd.Flags().GetString(utils.RESOURCE_KIND_KEY)

		err := controller.FindStaleResource(resourceKind, resourceGroup, resourceVersion, resourceName)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}

	},
}

var deleteCommand = &cobra.Command{
	Use:   "delete",
	Short: "Delete stale resource",
	Run: func(cmd *cobra.Command, args []string) {
		resourceName, _ := cmd.Flags().GetString(utils.RESOURCE_NAME_KEY)
		resourceVersion, _ := cmd.Flags().GetString(utils.RESOURCE_VERSION_KEY)
		resourceGroup, _ := cmd.Flags().GetString(utils.RESOURCE_GROUP_KEY)
		resourceKind, _ := cmd.Flags().GetString(utils.RESOURCE_KIND_KEY)

		fmt.Print(resourceName, resourceKind, resourceVersion, resourceGroup)
	},
}

func Execute() {
	var err error
	rootCmd := &cobra.Command{Use: "scour"}
	rootCmd.Version = utils.SCOUR_VERSION

	// add flags
	rootCmd.PersistentFlags().StringP(utils.RESOURCE_KIND_KEY, "k", "", "Resource Kind")
	err = rootCmd.MarkPersistentFlagRequired(utils.RESOURCE_KIND_KEY)
	if err != nil {
		panic(err)
	}

	rootCmd.PersistentFlags().StringArrayP("conditions", "c", []string{}, "Conditions")
	rootCmd.PersistentFlags().StringP(utils.RESOURCE_GROUP_KEY, "g", "", "Resource Group")
	rootCmd.PersistentFlags().StringP(utils.RESOURCE_VERSION_KEY, "v", "", "Resource Version")
	rootCmd.PersistentFlags().StringP(utils.RESOURCE_NAME_KEY, "n", "", "Resource Name")
	rootCmd.PersistentFlags().StringP(utils.RESOURCE_NAMESPACE_KEY, "s", "", "Resource Namespace")

	rootCmd.AddCommand(findCommand)

	if err = rootCmd.Execute(); err != nil {
		panic(err)
	}

}
