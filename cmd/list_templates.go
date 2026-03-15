package cmd

import (
	"fmt"
	tplfs "tfauto/templates"

	"github.com/spf13/cobra"
)

var list_templatesCmd = &cobra.Command{
	Use:   "templates",
	Short: "List available Terraform templates",
	Long:  "List all available built-in Terraform templates embedded in the tfauto binary.",
	RunE: func(cmd *cobra.Command, args []string) error {
		names, err := tplfs.List()
		if err != nil {
			return fmt.Errorf("list templates: %w", err)
		}

		fmt.Printf("Available templates:\n")
		for _, name := range names {
			fmt.Printf("  %s\n", name)
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
