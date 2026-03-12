package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/spf13/cobra"
)

// templatesCmd lists all available templates
var list_templatesCmd = &cobra.Command{
	Use:   "templates",
	Short: "List available Terraform templates",
	Long:  "Lists all available Terraform templates found in the ./templates directory.",
	RunE: func(cmd *cobra.Command, args []string) error {
		root := "templates"

		entries, err := os.ReadDir(root)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Println("No templates directory found at", root)
				fmt.Println("Create templates in ./templates/<name> to use tfauto init.")
				return nil
			}
			return fmt.Errorf("failed to read templates directory %q: %w", root, err)
		}

		var names []string
		for _, e := range entries {
			if e.IsDir() {
				names = append(names, e.Name())
			}
		}

		if len(names) == 0 {
			fmt.Println("No templates found in", root)
			fmt.Println("Add templates under ./templates/<name> to use with --template.")
			return nil
		}

		sort.Strings(names)

		fmt.Printf("Available templates:\n")
		for _, name := range names {
			fmt.Printf("  %-16s (folder: %s)\n", name, filepath.Join(root, name))
		}

		fmt.Println()
		fmt.Println("Use:")
		fmt.Println("  tfauto template <name>                    # show details and description")
		fmt.Println("  tfauto init --template <name> --target ./my-project")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(list_templatesCmd)
}
