package list

import (
	"bytes"
	"fmt"

	"github.com/maxthom/artifact-store/store"
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
	Run: func(cmd *cobra.Command, args []string) {
		var b bytes.Buffer
		e := yaml.NewEncoder(&b)
		e.SetIndent(2)

		store.InitializeStore()
		if len(args) > 2 {
			if args[2] == "." {
				if s, ok := store.ListFiles(args[0], args[1], withBundle); ok {
					if s, ok := store.ListFileContent(args[0], args[1], s...); ok {
						for k, v := range s {
							if len(s) > 1 {
								fmt.Println("--- " + k + " ------------")
							}
							fmt.Println(string(v))
						}
					} else {
						fmt.Println("no such bundle")
					}

				}
			} else if s, ok := store.ListFileContent(args[0], args[1], args[2:]...); ok {
				for k, v := range s {
					if len(s) > 1 {
						fmt.Println("--- " + k + " ------------")
					}
					fmt.Println(string(v))
				}
			} else {
				fmt.Println("no such bundle")
			}

		} else if len(args) == 2 {
			if s, ok := store.ListFiles(args[0], args[1], withBundle); ok {
				e.Encode(&s)
				fmt.Println(b.String())
			} else {
				fmt.Println("no such bundle")
			}
		} else if len(args) == 1 {
			s := store.ListBundles(args[0], "")
			for _, b := range s {
				fmt.Println("- " + b.Name + "/" + b.Version)
			}
		} else if len(args) == 0 {
			s := store.ListStore()
			for k, _ := range s {
				fmt.Println("- " + k)
			}
		}

	},
}
