package cmd

import (
	"fmt"
	"tfauto/internal/terraform"

	"github.com/spf13/cobra"
)

var validatePath string

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate Terraform configuration in a project directory",

	PreRunE: func(cmd *cobra.Command, args []string) error {
		if validatePath == "" {
			validatePath = "."
		}
		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		fmt.Println("Running terraform validate in", validatePath)

		if err := terraform.InitForValidation(ctx, validatePath); err != nil {
			return fmt.Errorf("terraform init for validate: %w", err)
		}

		if err := terraform.Validate(ctx, validatePath); err != nil {
			return fmt.Errorf("terraform validation: %w", err)
		}
		fmt.Println("validate: configuration is valid")

		return nil
	},
}

func init() {
	validateCmd.Flags().StringVar(&validatePath, "path", ".", "Path to Terraform project")
	rootCmd.AddCommand(validateCmd)
}
