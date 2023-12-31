package list

import (
	"bytes"
	"fmt"

	"github.com/maxthom/artifact-store/store"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var ()

func init() {

}

var bundleCmd = &cobra.Command{
	Use:   "bundles",
	Short: "List bundles in store",
	Args:  cobra.MaximumNArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		var b bytes.Buffer
		e := yaml.NewEncoder(&b)
		e.SetIndent(2)

		if len(args) == 0 {
			s := store.ListStore()
			e.Encode(&s)
		} else if len(args) == 1 {
			s := store.ListBundles(args[0], "", "")
			e.Encode(&s)
		} else if len(args) == 2 {
			s := store.ListBundles(args[0], args[1], "")
			e.Encode(&s)
		} else if len(args) == 3 {
			s := store.ListBundles(args[0], args[1], args[2])
			e.Encode(&s)
		}
		fmt.Println(b.String())
	},
}
