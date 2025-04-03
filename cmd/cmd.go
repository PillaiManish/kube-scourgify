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
		//TODO: resourceName, _ := cmd.Flags().GetString(utils.RESOURCE_NAME_KEY)
		resourceVersion, _ := cmd.Flags().GetString(utils.RESOURCE_VERSION_KEY)
		resourceGroup, _ := cmd.Flags().GetString(utils.RESOURCE_GROUP_KEY)
		resourceKind, _ := cmd.Flags().GetString(utils.RESOURCE_KIND_KEY)
		conditionsFilepath, _ := cmd.Flags().GetString(utils.CONDITIONS_FILEPATH)

		err := controller.FindStaleResource(resourceKind, resourceGroup, resourceVersion, conditionsFilepath)
		if err != nil {
			panic(err)
		}

	},
}

var deleteCommand = &cobra.Command{
	Use:   "delete",
	Short: "Delete stale resource",
	Run: func(cmd *cobra.Command, args []string) {
		//TODO: resourceName, _ := cmd.Flags().GetString(utils.RESOURCE_NAME_KEY)
		//resourceVersion, _ := cmd.Flags().GetString(utils.RESOURCE_VERSION_KEY)
		//resourceGroup, _ := cmd.Flags().GetString(utils.RESOURCE_GROUP_KEY)
		//resourceKind, _ := cmd.Flags().GetString(utils.RESOURCE_KIND_KEY)
		resourceCRName, _ := cmd.Flags().GetString(utils.RESOURCE_CR_NAME)
		fmt.Print(resourceCRName)
	},
}

func Execute() {
	var err error
	rootCmd := &cobra.Command{Use: "scour", CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true}}
	rootCmd.Version = utils.SCOUR_VERSION
	
	findCommand.PersistentFlags().StringP(utils.CONDITIONS_FILEPATH, "c", "", "absolute filepath to conditions.json")
	findCommand.PersistentFlags().StringP(utils.RESOURCE_GROUP_KEY, "g", "", "Resource Group")
	findCommand.PersistentFlags().StringP(utils.RESOURCE_VERSION_KEY, "v", "", "Resource Version")
	findCommand.PersistentFlags().StringP(utils.RESOURCE_KIND_KEY, "n", "", "Resource Kind")
	findCommand.PersistentFlags().StringP(utils.RESOURCE_NAMESPACE_KEY, "s", "", "Resource Namespace")

	err = findCommand.MarkPersistentFlagRequired(utils.CONDITIONS_FILEPATH)
	if err != nil {
		panic(err)
	}

	rootCmd.AddCommand(findCommand, deleteCommand)

	if err = rootCmd.Execute(); err != nil {
		panic(err)
	}

}
