package cmd

import (
	"fmt"
	"strings"

	tplfs "github.com/khushaltarsariya/tfauto/templates"

	"github.com/spf13/cobra"
)

// templateCmd shows detailed information about a single template.
// Usage: tfauto template <name>
var use_templateCmd = &cobra.Command{
	Use:     "template [name]",
	Aliases: []string{"show", "inspect"},
	Short:   "Show metadata and files for a built-in Terraform template",
	Long: `Show detailed information about a Terraform template.

It reads template.yaml or template.json when present, then falls back to
legacy metadata inferred from DESCRIPTION.md and provider files.

This keeps older templates working while allowing new templates to ship a
manifest with richer metadata.`,
	Example: `  tfauto template aws-basic
  tfauto template aws-basic --json`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		if !tplfs.Exists(name) {
			return fmt.Errorf("tfauto template: template %q not found", name)
		}

		metadata, err := tplfs.Metadata(name)
		if err != nil {
			return fmt.Errorf("tfauto template: load template metadata: %w", err)
		}
		files, err := tplfs.Files(name)
		if err != nil {
			return fmt.Errorf("tfauto template: list template files: %w", err)
		}

		if jsonRequested(cmd) {
			return writeJSON(cmd.OutOrStdout(), map[string]any{
				"command":  "template",
				"template": metadata,
				"files":    files,
			})
		}

		fmt.Printf("tfauto: template %s\n", metadata.Name)
		fmt.Println()
		fmt.Println("Metadata:")
		fmt.Printf("  Name: %s\n", displayField(metadata.Name))
		fmt.Printf("  Version: %s\n", displayField(metadata.Version))
		fmt.Printf("  Author: %s\n", displayField(metadata.Author))
		fmt.Printf("  Category: %s\n", displayField(metadata.Category))
		fmt.Printf("  Cloud provider: %s\n", displayField(metadata.CloudProvider))
		fmt.Printf("  Estimated monthly cost: %s\n", displayField(metadata.EstimatedMonthlyCost))
		fmt.Printf("  Required Terraform version: %s\n", displayField(metadata.RequiredTerraformVersion))
		fmt.Printf("  Required providers: %s\n", joinDisplay(metadata.RequiredProviders))
		fmt.Printf("  Tags: %s\n", joinDisplay(metadata.Tags))
		fmt.Printf("  Metadata source: %s\n", displayField(metadata.MetadataSource))
		if metadata.Description != "" {
			fmt.Println()
			fmt.Println("Description:")
			fmt.Println(strings.TrimSpace(metadata.Description))
		}
		fmt.Println()
		fmt.Println("Files:")
		for _, file := range files {
			fmt.Printf("  - %s\n", file)
		}

		fmt.Println()
		fmt.Println("Next step:")
		fmt.Printf("  tfauto init --template %s --target ./my-project\n", name)

		return nil
	},
}

func init() {
	use_templateCmd.Flags().Bool("json", false, "Output template details as JSON")
	rootCmd.AddCommand(use_templateCmd)
}

func joinDisplay(values []string) string {
	if len(values) == 0 {
		return "-"
	}
	return strings.Join(values, ", ")
}
