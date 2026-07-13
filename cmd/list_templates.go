package cmd

import (
	"fmt"
	"strings"
	"text/tabwriter"

	tplfs "github.com/khushaltarsariya/tfauto/templates"

	"github.com/spf13/cobra"
)

var list_templatesCmd = &cobra.Command{
	Use:     "templates",
	Aliases: []string{"list", "ls"},
	Short:   "List built-in Terraform templates",
	Long:    "List all available built-in Terraform templates and their metadata.",
	Example: `  tfauto templates
  tfauto templates --json`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		templates, err := tplfs.ListMetadata()
		if err != nil {
			return fmt.Errorf("tfauto templates: %w", err)
		}

		if jsonRequested(cmd) {
			return writeJSON(cmd.OutOrStdout(), map[string]any{
				"command":   "templates",
				"templates": templates,
			})
		}

		fmt.Println("tfauto: templates")
		fmt.Println()
		w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tVERSION\tPROVIDER\tTF VER\tCATEGORY\tCOST\tDESCRIPTION")
		for _, item := range templates {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
				displayField(item.Name),
				displayField(item.Version),
				displayField(item.CloudProvider),
				displayField(item.RequiredTerraformVersion),
				displayField(item.Category),
				displayField(item.EstimatedMonthlyCost),
				shortenField(item.Description, 56),
			)
		}
		_ = w.Flush()
		fmt.Println()
		fmt.Println("Next steps:")
		fmt.Println("  tfauto template <name>")
		fmt.Println("  tfauto init --template <name> --target ./my-project")

		return nil
	},
}

func init() {
	list_templatesCmd.Flags().Bool("json", false, "Output template list as JSON")
	rootCmd.AddCommand(list_templatesCmd)
}

func displayField(value string) string {
	if value == "" {
		return "-"
	}
	return value
}

func shortenField(value string, limit int) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "-"
	}
	if len(value) <= limit {
		return value
	}
	if limit <= 1 {
		return value[:1]
	}
	return value[:limit-1] + "..."
}
