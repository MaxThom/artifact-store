package list

import (
	"bytes"
	"fmt"

	"github.com/maxthom/artifact-store/services"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	withBundle bool
)

func init() {
	filesCmd.SetUsageFunc(func(c *cobra.Command) error {
		fmt.Print("Usage: artifact-store list files <bundle> <version> <file> [flags]\n\n")
		fmt.Print("Flags:\n")
		fmt.Print("  -h, --help      help for files\n")
		fmt.Print("  --with-bundle   include bundle files in output\n")
		return nil
	})

	filesCmd.Flags().BoolVar(&withBundle, "with-bundle", false, "include bundle files in output")
}

var filesCmd = &cobra.Command{
	Use:   "files",
	Short: "List files in a bundle",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		var b bytes.Buffer
		e := yaml.NewEncoder(&b)
		e.SetIndent(2)

		services.InitializeStore()
		if len(args) > 2 {
			if args[2] == "." {
				if s, ok := services.ListFiles(args[0], args[1], withBundle); ok {
					if s, ok := services.ListFileContent(args[0], args[1], s...); ok {
						for k, v := range s {
							fmt.Println("--- " + k + " ------------")
							fmt.Println(string(v))
						}
					} else {
						fmt.Println("no such bundle")
					}

				}
			} else if s, ok := services.ListFileContent(args[0], args[1], args[2:]...); ok {
				for k, v := range s {
					fmt.Println("--- " + k + " ------------")
					fmt.Println(string(v))
				}
			} else {
				fmt.Println("no such bundle")
			}

		} else if len(args) == 2 {
			if s, ok := services.ListFiles(args[0], args[1], withBundle); ok {
				e.Encode(&s)
				fmt.Println(b.String())
			} else {
				fmt.Println("no such bundle")
			}

		}
	},
}