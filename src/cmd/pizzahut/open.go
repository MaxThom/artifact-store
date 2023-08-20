package pizzahut

import (
	"github.com/maxthom/artifact-store/pizzahut"
	"github.com/spf13/cobra"
)

var ()

func init() {

}

var openCmd = &cobra.Command{
	Use:   "open",
	Short: "Open store for the day",
	Run: func(cmd *cobra.Command, args []string) {
		pizzahut.OpenStore()
	},
}
