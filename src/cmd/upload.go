package cmd

import (
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

		services.InitializeStore()

		services.UploadBundle(args[0], path)

	},
}
