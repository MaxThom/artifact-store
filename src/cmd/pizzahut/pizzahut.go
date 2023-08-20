package pizzahut

import (
	"github.com/spf13/cobra"
)

var ()

func init() {

	PizzahutCmd.AddCommand(openCmd)
	PizzahutCmd.AddCommand(closeCmd)
}

var PizzahutCmd = &cobra.Command{
	Use:   "pizzahut",
	Short: "Manage pizzahut store",
}
