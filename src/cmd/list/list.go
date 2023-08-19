package list

import (
	"github.com/spf13/cobra"
)

var ()

func init() {

	ListCmd.AddCommand(bundleCmd)
	ListCmd.AddCommand(filesCmd)
}

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List bundles and files in store",
}
