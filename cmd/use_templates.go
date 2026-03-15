package cmd

import (
	"fmt"
	"strings"
	tplfs "tfauto/templates"

	"github.com/spf13/cobra"
)

// templateCmd shows detailed information about a single template.
// Usage: tfauto template <name>
var use_templateCmd = &cobra.Command{
	Use:   "template [name]",
	Short: "Show detailed information about a Terraform template",
	Long: `Show detailed information about a Terraform template.

It reads DESCRIPTION.md from the template folder (if present) and
lists all files included in the template.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		if !tplfs.Exists(name) {
			return fmt.Errorf("template %q not found", name)
		}

		fmt.Printf("Template: %s\n", name)
		fmt.Printf("Source:   embedded in tfauto binary\n\n")

		if description, err := tplfs.Description(name); err == nil && description != "" {
			fmt.Println("Description:")
			fmt.Println(strings.TrimSpace(description))
			fmt.Println()
		} else if err == nil {
			fmt.Println("Description:")
			fmt.Println("(no DESCRIPTION.md found for this template)")
			fmt.Println()
		} else {
			return fmt.Errorf("read template description: %w", err)
		}

		fmt.Println("Files in this template:")
		files, err := tplfs.Files(name)
		if err != nil {
			return fmt.Errorf("list template files: %w", err)
		}
		for _, file := range files {
			fmt.Printf("  - %s\n", file)
		}

		fmt.Println()
		fmt.Println("You can use this template with:")
		fmt.Println("  tfauto init --template", name, "--target ./my-project")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(use_templateCmd)
}
