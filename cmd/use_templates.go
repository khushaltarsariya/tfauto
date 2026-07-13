package cmd

import (
	"fmt"
	"io"
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
				"ok":       true,
				"template": metadata,
				"files":    files,
			})
		}

		renderTemplate(cmd.OutOrStdout(), metadata, files, name)

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

func renderTemplate(w io.Writer, metadata tplfs.TemplateMetadata, files []string, name string) {
	fmt.Fprintf(w, "tfauto: template %s\n", metadata.Name)
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Metadata:")
	fmt.Fprintf(w, "  Name: %s\n", displayField(metadata.Name))
	fmt.Fprintf(w, "  Version: %s\n", displayField(metadata.Version))
	fmt.Fprintf(w, "  Author: %s\n", displayField(metadata.Author))
	fmt.Fprintf(w, "  Category: %s\n", displayField(metadata.Category))
	fmt.Fprintf(w, "  Cloud provider: %s\n", displayField(metadata.CloudProvider))
	fmt.Fprintf(w, "  Estimated monthly cost: %s\n", displayField(metadata.EstimatedMonthlyCost))
	fmt.Fprintf(w, "  Required Terraform version: %s\n", displayField(metadata.RequiredTerraformVersion))
	fmt.Fprintf(w, "  Required providers: %s\n", joinDisplay(metadata.RequiredProviders))
	fmt.Fprintf(w, "  Tags: %s\n", joinDisplay(metadata.Tags))
	fmt.Fprintf(w, "  Metadata source: %s\n", displayField(metadata.MetadataSource))
	if metadata.Description != "" {
		fmt.Fprintln(w)
		fmt.Fprintln(w, "Description:")
		fmt.Fprintln(w, strings.TrimSpace(metadata.Description))
	}
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Files:")
	for _, file := range files {
		fmt.Fprintf(w, "  - %s\n", file)
	}

	fmt.Fprintln(w)
	fmt.Fprintln(w, "Next step:")
	fmt.Fprintf(w, "  tfauto init --template %s --target ./my-project\n", name)
}
