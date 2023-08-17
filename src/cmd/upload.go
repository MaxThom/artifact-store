package cmd

import (
	"fmt"

	"github.com/maxthom/artifact-store/services"
	"github.com/spf13/cobra"
)

var (
	path string
)

func init() {
	uploadCmd.Flags().StringVarP(&path, "storedPath", "p", ".", "relative path from store root")
	rootCmd.AddCommand(uploadCmd)
}

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload bundle to server",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(args)
		fmt.Println(path)

		services.InitializeStore()

	},
}
