package cmd

import (
	"fmt"
	"tfauto/internal/terraform"

	"github.com/spf13/cobra"
)

var validatePath string

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "to validate the tf files",

	PreRunE: func(cmd *cobra.Command, args []string) error {
		if validatePath == "" {
			validatePath = "."
		}
		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		fmt.Println("run terraform validation command", validatePath)

		if err := terraform.Validate(ctx, validatePath); err != nil {
			return fmt.Errorf("terraform validation: %w", err)
		}
		fmt.Println("validate: configuration is valid")

		return nil
	},
}

func init() {
	validateCmd.Flags().StringVar(&validatePath, "path", ".", "path to terraform project")
	rootCmd.AddCommand(validateCmd)
}
