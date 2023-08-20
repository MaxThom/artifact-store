package pizzahut

import (
	"github.com/maxthom/artifact-store/pizzahut"
	"github.com/spf13/cobra"
)

var ()

func init() {

}

var closeCmd = &cobra.Command{
	Use:   "close",
	Short: "Close store for the day",
	Run: func(cmd *cobra.Command, args []string) {
		pizzahut.CloseStore()
	},
}
