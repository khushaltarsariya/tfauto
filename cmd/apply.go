package cmd

import (
	"fmt"
	"tfauto/internal/terraform"

	"github.com/spf13/cobra"
)

var applyPath string

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Run terraform apply in a project directory",

	PreRunE: func(cmd *cobra.Command, args []string) error {
		if applyPath == "" {
			applyPath = "."
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		fmt.Println("Running terraform apply in", applyPath)

		if err := terraform.Init(ctx, applyPath); err != nil {
			return fmt.Errorf("terraform init failed: %w", err)
		}

		if err := terraform.Apply(ctx, applyPath); err != nil {
			return fmt.Errorf("terraform apply failed: %w", err)
		}
		return nil
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		fmt.Println("Apply completed")
	},
}

func init() {
	applyCmd.Flags().StringVar(&applyPath, "path", ".", "Path to Terraform project")
	rootCmd.AddCommand(applyCmd)
}
