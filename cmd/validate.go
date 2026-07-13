package cmd

import (
	"fmt"

	"github.com/khushaltarsariya/tfauto/internal/terraform"

	"github.com/spf13/cobra"
)

var validatePath string
var validateJSON bool

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate Terraform configuration for a project directory",
	Long: `Validate Terraform configuration in a project directory.

Examples:
  tfauto validate --path ./app
  tfauto validate --path ./app --json`,
	Example: `  tfauto validate --path ./app`,
	Args:    cobra.NoArgs,

	PreRunE: func(cmd *cobra.Command, args []string) error {
		if validatePath == "" {
			validatePath = "."
		}
		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		if err := terraform.InitForValidation(ctx, validatePath); err != nil {
			return fmt.Errorf("tfauto validate: terraform init for validation failed: %w", err)
		}

		if err := terraform.Validate(ctx, validatePath); err != nil {
			return fmt.Errorf("tfauto validate: terraform validation failed: %w", err)
		}
		if validateJSON || jsonRequested(cmd) {
			return writeJSON(cmd.OutOrStdout(), map[string]any{
				"command": "validate",
				"ok":      true,
				"path":    validatePath,
			})
		}
		fmt.Println(tfautoMessage("validate", "completed successfully"))

		return nil
	},
}

func init() {
	validateCmd.Flags().StringVar(&validatePath, "path", ".", "Path to Terraform project")
	validateCmd.Flags().BoolVar(&validateJSON, "json", false, "Output validation details as JSON")
	rootCmd.AddCommand(validateCmd)
}
