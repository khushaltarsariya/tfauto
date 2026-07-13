package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/khushaltarsariya/tfauto/internal/terraform"

	"github.com/spf13/cobra"
)

var applyPath string
var applyPlanFile string
var applyYes bool
var applyRequirePlan bool
var applyJSON bool

var applyCmd = &cobra.Command{
	Use:     "apply",
	Aliases: []string{"deploy"},
	Short:   "Apply Terraform changes in a project directory",
	Long: `Run terraform apply in a project directory.

Examples:
  tfauto apply --path ./app --yes
  tfauto apply --path ./app --plan tfplan
  tfauto apply --path ./app --json`,
	Example: `  tfauto apply --path ./app --yes
  tfauto apply --path ./app --plan tfplan`,
	Args: cobra.NoArgs,

	PreRunE: func(cmd *cobra.Command, args []string) error {
		if applyPath == "" {
			applyPath = "."
		}

		configResult, err := loadConfigForPath(applyPath)
		if err != nil {
			return tfautoError("apply", err)
		}
		if configResult.Found && configResult.Config.Terraform.RequirePlanFile {
			applyRequirePlan = true
		}

		if applyRequirePlan && applyPlanFile == "" {
			return fmt.Errorf("tfauto apply: saved plan required; pass --plan because --require-plan or project policy is enabled")
		}

		if applyPlanFile != "" {
			planPath := resolvePathFromProject(applyPath, applyPlanFile)
			if _, err := os.Stat(planPath); err != nil {
				return fmt.Errorf("tfauto apply: plan file %q does not exist", applyPlanFile)
			}
		}

		if applyPlanFile == "" && !applyYes {
			return fmt.Errorf("tfauto apply: non-interactive apply requires --yes or --plan")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		if err := terraform.Init(ctx, applyPath); err != nil {
			return fmt.Errorf("tfauto apply: terraform init failed: %w", err)
		}

		if applyPlanFile != "" {
			if err := terraform.ApplyPlan(ctx, applyPath, applyPlanFile); err != nil {
				return fmt.Errorf("tfauto apply: terraform apply saved plan failed: %w", err)
			}
			if applyJSON || jsonRequested(cmd) {
				return writeJSON(cmd.OutOrStdout(), map[string]any{
					"command":    "apply",
					"ok":         true,
					"path":       applyPath,
					"plan":       applyPlanFile,
					"saved_plan": true,
				})
			}
			fmt.Println(tfautoMessage("apply", "saved plan applied successfully"))
			return nil
		}

		if err := terraform.Apply(ctx, applyPath); err != nil {
			return fmt.Errorf("tfauto apply: terraform apply failed: %w", err)
		}
		if applyJSON || jsonRequested(cmd) {
			return writeJSON(cmd.OutOrStdout(), map[string]any{
				"command": "apply",
				"ok":      true,
				"path":    applyPath,
			})
		}
		fmt.Println(tfautoMessage("apply", "completed successfully"))
		return nil
	},
}

func init() {
	applyCmd.Flags().StringVar(&applyPath, "path", ".", "Path to Terraform project")
	applyCmd.Flags().StringVar(&applyPlanFile, "plan", "", "Apply a previously saved plan file")
	applyCmd.Flags().BoolVar(&applyYes, "yes", false, "Confirm non-interactive apply when no saved plan file is provided")
	applyCmd.Flags().BoolVar(&applyRequirePlan, "require-plan", false, "Require --plan for safer CI usage")
	applyCmd.Flags().BoolVar(&applyJSON, "json", false, "Output apply details as JSON")
	rootCmd.AddCommand(applyCmd)
}

func resolvePathFromProject(projectPath string, value string) string {
	if filepath.IsAbs(value) {
		return value
	}
	return filepath.Join(projectPath, value)
}
