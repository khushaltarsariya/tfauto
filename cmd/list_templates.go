package cmd

import (
	"fmt"
	tplfs "tfauto/templates"

	"github.com/spf13/cobra"
)

var list_templatesCmd = &cobra.Command{
	Use:     "templates",
	Aliases: []string{"list", "ls"},
	Short:   "List built-in Terraform templates",
	Long:    "List all available built-in Terraform templates embedded in the tfauto binary.",
	Example: `  tfauto templates
  tfauto templates --json`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		names, err := tplfs.List()
		if err != nil {
			return fmt.Errorf("tfauto templates: %w", err)
		}

		if jsonRequested(cmd) {
			return writeJSON(cmd.OutOrStdout(), map[string]any{
				"command":   "templates",
				"templates": names,
			})
		}

		fmt.Println("tfauto: templates")
		fmt.Println()
		fmt.Println("Available templates:")
		for _, name := range names {
			fmt.Printf("  - %s\n", name)
		}
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
