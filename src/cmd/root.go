package cmd

import (
	"fmt"
	"os"

	"github.com/maxthom/artifact-store/cmd/list"
	"github.com/maxthom/artifact-store/cmd/pizzahut"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "artifact-store",
	Short: "Manage a set of artifacts on disk or s3",
}

func init() {
	rootCmd.AddCommand(list.ListCmd)
	rootCmd.AddCommand(pizzahut.PizzahutCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
