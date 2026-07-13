package cmd

import (
	"fmt"

	"github.com/khushaltarsariya/tfauto/internal/terraform"

	"github.com/spf13/cobra"
)

var pathFlag string
var planOut string
var planDetailedExitCode bool
var planJSON bool

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Generate a Terraform plan",
	Long: `Run terraform plan in a project directory.

Examples:
  tfauto plan --path ./app
  tfauto plan --path ./app --out tfplan --detailed-exitcode
  tfauto plan --path ./app --json`,
	Example: `  tfauto plan --path ./app
  tfauto plan --path ./app --out tfplan --detailed-exitcode
  tfauto plan --path ./app --json`,
	Args: cobra.NoArgs,

	PreRunE: func(cmd *cobra.Command, args []string) error {
		if pathFlag == "" {
			pathFlag = "."
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		if err := terraform.Init(ctx, pathFlag); err != nil {
			return fmt.Errorf("tfauto plan: terraform init: %w", err)

		}

		if planOut != "" {
			result, err := terraform.PlanOut(ctx, pathFlag, planOut, planDetailedExitCode)
			if err != nil {
				return fmt.Errorf("tfauto plan: terraform plan with -out: %w", err)
			}
			if result.HasChanges {
				if planJSON || jsonRequested(cmd) {
					return writeJSON(cmd.OutOrStdout(), map[string]any{
						"command":     "plan",
						"ok":          true,
						"path":        pathFlag,
						"out":         planOut,
						"has_changes": true,
						"exit_code":   2,
					})
				}
				fmt.Println(tfautoMessage("plan", "terraform plan detected pending changes"))
				fmt.Println(tfautoMessage("plan", "wrote plan to %s", planOut))
				return ExitError{Code: 2}
			}
			if planJSON || jsonRequested(cmd) {
				return writeJSON(cmd.OutOrStdout(), map[string]any{
					"command":     "plan",
					"ok":          true,
					"path":        pathFlag,
					"out":         planOut,
					"has_changes": false,
					"exit_code":   0,
				})
			}
			fmt.Println(tfautoMessage("plan", "wrote plan to %s", planOut))
			fmt.Println(tfautoMessage("plan", "completed successfully"))
			return nil
		}
		result, err := terraform.Plan(ctx, pathFlag, planDetailedExitCode)
		if err != nil {
			return fmt.Errorf("tfauto plan: terraform plan: %w", err)
		}
		if result.HasChanges {
			if planJSON || jsonRequested(cmd) {
				return writeJSON(cmd.OutOrStdout(), map[string]any{
					"command":     "plan",
					"ok":          true,
					"path":        pathFlag,
					"has_changes": true,
					"exit_code":   2,
				})
			}
			fmt.Println(tfautoMessage("plan", "terraform plan detected pending changes"))
			return ExitError{Code: 2}
		}
		if planJSON || jsonRequested(cmd) {
			return writeJSON(cmd.OutOrStdout(), map[string]any{
				"command":     "plan",
				"ok":          true,
				"path":        pathFlag,
				"has_changes": false,
				"exit_code":   0,
			})
		}
		fmt.Println(tfautoMessage("plan", "completed successfully"))
		return nil
	},
}

func init() {
	planCmd.Flags().StringVar(&pathFlag, "path", ".", "Path to Terraform project")
	planCmd.Flags().StringVar(&planOut, "out", "", "Write the plan to a file")
	planCmd.Flags().BoolVar(&planDetailedExitCode, "detailed-exitcode", false, "Return exit code 2 when the plan contains changes")
	planCmd.Flags().BoolVar(&planJSON, "json", false, "Output plan details as JSON")
	rootCmd.AddCommand(planCmd)
}
