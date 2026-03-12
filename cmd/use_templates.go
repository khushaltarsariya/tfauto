package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

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
		root := "templates"
		templateDir := filepath.Join(root, name)

		info, err := os.Stat(templateDir)
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("template %q not found in %s", name, root)
			}
			return fmt.Errorf("unable to stat template directory %q: %w", templateDir, err)
		}
		if !info.IsDir() {
			return fmt.Errorf("template %q is not a directory", templateDir)
		}

		fmt.Printf("Template: %s\n", name)
		fmt.Printf("Path:     %s\n\n", templateDir)

		// Read DESCRIPTION.md if it exists
		descPath := filepath.Join(templateDir, "DESCRIPTION.md")
		if data, err := os.ReadFile(descPath); err == nil {
			fmt.Println("Description:")
			fmt.Println(strings.TrimSpace(string(data)))
			fmt.Println()
		} else if os.IsNotExist(err) {
			fmt.Println("Description:")
			fmt.Println("(no DESCRIPTION.md found for this template)")
			fmt.Println()
		} else {
			return fmt.Errorf("failed to read DESCRIPTION.md: %w", err)
		}

		// List files inside template directory
		fmt.Println("Files in this template:")
		err = filepath.WalkDir(templateDir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}

			rel, err := filepath.Rel(templateDir, path)
			if err != nil {
				return err
			}

			// Skip DESCRIPTION.md in the files list if you want
			if rel == "DESCRIPTION.md" {
				return nil
			}

			fmt.Printf("  - %s\n", rel)
			return nil
		})
		if err != nil {
			return fmt.Errorf("failed to walk template files: %w", err)
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
