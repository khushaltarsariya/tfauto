package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"tfauto/internal/config"
	"tfauto/internal/terraform"

	"github.com/spf13/cobra"
)

var applyPath string
var applyPlanFile string
var applyYes bool
var applyRequirePlan bool

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Run terraform apply in a project directory",

	PreRunE: func(cmd *cobra.Command, args []string) error {
		if applyPath == "" {
			applyPath = "."
		}

		configResult, err := config.LoadForPath(applyPath)
		if err != nil {
			return err
		}
		if configResult.Found && configResult.Config.Terraform.RequirePlanFile {
			applyRequirePlan = true
			fmt.Println("Project policy requires applying a saved plan file:", configResult.Path)
		}

		if applyRequirePlan && applyPlanFile == "" {
			return fmt.Errorf("apply requires --plan because --require-plan or project policy is enabled")
		}

		if applyPlanFile != "" {
			planPath := resolvePathFromProject(applyPath, applyPlanFile)
			if _, err := os.Stat(planPath); err != nil {
				return fmt.Errorf("plan file %q does not exist", applyPlanFile)
			}
		}

		if applyPlanFile == "" && !applyYes {
			return fmt.Errorf("apply without a saved plan requires --yes")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		fmt.Println("Running terraform apply in", applyPath)

		if err := terraform.Init(ctx, applyPath); err != nil {
			return fmt.Errorf("terraform init failed: %w", err)
		}

		if applyPlanFile != "" {
			if err := terraform.ApplyPlan(ctx, applyPath, applyPlanFile); err != nil {
				return fmt.Errorf("terraform apply saved plan failed: %w", err)
			}
			return nil
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
	applyCmd.Flags().StringVar(&applyPlanFile, "plan", "", "Apply a previously saved plan file")
	applyCmd.Flags().BoolVar(&applyYes, "yes", false, "Confirm non-interactive apply when no saved plan file is provided")
	applyCmd.Flags().BoolVar(&applyRequirePlan, "require-plan", false, "Require --plan for safer CI usage")
	rootCmd.AddCommand(applyCmd)
}

func resolvePathFromProject(projectPath string, value string) string {
	if filepath.IsAbs(value) {
		return value
	}
	return filepath.Join(projectPath, value)
}
