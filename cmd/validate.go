package cmd

import (
	"fmt"
	"tfauto/internal/terraform"

	"github.com/spf13/cobra"
)

var validatePath string

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate Terraform configuration for a project directory",
	Long: `Validate Terraform configuration in a project directory.

Examples:
  tfauto validate --path ./app`,
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
		fmt.Println(tfautoMessage("validate", "completed successfully"))

		return nil
	},
}

func init() {
	validateCmd.Flags().StringVar(&validatePath, "path", ".", "Path to Terraform project")
	rootCmd.AddCommand(validateCmd)
}
